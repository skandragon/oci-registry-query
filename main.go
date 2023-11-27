package main

import (
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func main() {
	ref, err := name.ParseReference(os.Args[1])
	check(err)

	raw, err := remote.Get(ref)
	check(err)
	switch raw.MediaType {
	case types.OCIImageIndex, types.DockerManifestList:
		index, err := remote.Index(ref)
		check(err)

		im, err := index.IndexManifest()
		check(err)

		for i, m := range im.Manifests {
			log.Printf("%d MANIFEST: size=%d digest=%s mediaType=%s", i, m.Size, m.Digest, m.MediaType)
			log.Printf("  %s/%s", m.Platform.OS, m.Platform.Architecture)
		}
	case types.OCIManifestSchema1, types.DockerManifestSchema1, types.DockerManifestSchema2:
		img, err := remote.Image(ref)
		check(err)

		m, err := img.Manifest()
		check(err)

		for i, l := range m.Layers {
			log.Printf("%d LAYER: size=%d digest=%s mediaType=%s", i, l.Size, l.Digest, l.MediaType)
		}

	default:
		log.Printf("Found media type %s", raw.MediaType)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

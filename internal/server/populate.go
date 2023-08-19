package server

import (
	"filesync/internal/utils"
	"filesync/pkg/delta"
	"log"

	"github.com/cespare/xxhash/v2"
)

func populate(root string, change utils.Change) utils.Change {
	if change.IsDir {
		return change
	}
	if change.Type == utils.Create {
		emptySignature := delta.Signature{
			Hash:        xxhash.Sum64([]byte{}),
			BlockSize:   0,
			BlockHashes: map[uint64]delta.BlockInfo{},
		}
		patch, err := delta.CreatePatch(emptySignature, root+"/"+change.ToPath)
		if err != nil {
			log.Fatalln(err)
		}
		change.Patch = patch
	} else if change.Type == utils.Update {
		signature := delta.Signature{}

		err := delta.Decode[delta.Signature](change.FromSignature, &signature)
		if err != nil {
			log.Fatalln(err)
		}

		patch, err := delta.CreatePatch(signature, root+"/"+change.ToPath)
		if err != nil {
			log.Fatalln(err)
		}

		change.Patch = patch
	}
	return change
}

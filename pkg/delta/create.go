package delta

import (
	"os"

	"github.com/cespare/xxhash/v2"
)

// Reads the signature of the previous version of the file at newFilePath and the data of the file at newFilePath.
// Generates a patch to convert from the previous version to the current version and writes it to patchPath.
func CreatePatch(signature Signature, newFilePath string) ([]byte, error) {
	data, err := os.ReadFile(newFilePath)
	if err != nil {
		return nil, err
	}

	patch := Patch{}
	patch.FromHash = signature.Hash
	patch.ToHash = xxhash.Sum64(data)
	patch.Chunks = []Chunk{}

	start := 0
	for start < len(data) {
		end := start + signature.BlockSize
		if end > len(data) {
			end = len(data)
		}
		block := data[start:end]
		hash := xxhash.Sum64(block)
		info, exists := signature.BlockHashes[hash]
		if exists {
			chunk := Chunk{}
			chunk.IsBlock = true
			chunk.Block = info
			patch.Chunks = append(patch.Chunks, chunk)
			start += signature.BlockSize
		} else {
			if len(patch.Chunks) == 0 || patch.Chunks[len(patch.Chunks)-1].IsBlock {
				chunk := Chunk{}
				chunk.IsBlock = false
				chunk.Data = []byte{}
				patch.Chunks = append(patch.Chunks, chunk)
			}
			patch.Chunks[len(patch.Chunks)-1].Data = append(patch.Chunks[len(patch.Chunks)-1].Data, data[start])
			start++
		}
	}

	uncompressed, err := Encode(patch)
	if err != nil {
		return nil, err
	}

	compressed, err := compress(uncompressed)
	if err != nil {
		return nil, err
	}

	return compressed, nil
}

func CreatePatchFile(signature Signature, newFilePath string, patchPath string) error {
	data, err := os.ReadFile(newFilePath)
	if err != nil {
		return err
	}

	patch := Patch{}
	patch.FromHash = signature.Hash
	patch.ToHash = xxhash.Sum64(data)
	patch.Chunks = []Chunk{}

	start := 0
	for start < len(data) {
		end := start + signature.BlockSize
		if end > len(data) {
			end = len(data)
		}
		block := data[start:end]
		hash := xxhash.Sum64(block)
		info, exists := signature.BlockHashes[hash]
		if exists {
			chunk := Chunk{}
			chunk.IsBlock = true
			chunk.Block = info
			patch.Chunks = append(patch.Chunks, chunk)
			start += signature.BlockSize
		} else {
			if len(patch.Chunks) == 0 || patch.Chunks[len(patch.Chunks)-1].IsBlock {
				chunk := Chunk{}
				chunk.IsBlock = false
				chunk.Data = []byte{}
				patch.Chunks = append(patch.Chunks, chunk)
			}
			patch.Chunks[len(patch.Chunks)-1].Data = append(patch.Chunks[len(patch.Chunks)-1].Data, data[start])
			start++
		}
	}

	uncompressed, err := Encode(patch)
	if err != nil {
		return err
	}

	compressed, err := compress(uncompressed)
	if err != nil {
		return err
	}

	file, err := os.Create(patchPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(compressed)
	if err != nil {
		return err
	}

	return nil
}

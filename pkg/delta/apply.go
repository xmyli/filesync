package delta

import (
	"errors"
	"os"

	"github.com/cespare/xxhash/v2"
)

// Reads the patch file at patchPath and applies it to the file at inPath.
// The original file at inPath is unchanged and the patched file is written to outPath.
func ApplyPatch(patchBytes []byte, inPath string, outPath string) error {
	uncompressed, err := decompress(patchBytes)
	if err != nil {
		return err
	}

	patch := Patch{}

	err = Decode[Patch](uncompressed, &patch)
	if err != nil {
		return err
	}

	inBytes, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	outBytes := []byte{}

	if xxhash.Sum64(inBytes) != patch.FromHash {
		return errors.New("invalid input file")
	}

	for _, chunk := range patch.Chunks {
		if chunk.IsBlock {
			outBytes = append(outBytes, inBytes[chunk.Block.Start:chunk.Block.End]...)
		} else {
			outBytes = append(outBytes, chunk.Data...)
		}
	}

	if xxhash.Sum64(outBytes) != patch.ToHash {
		return errors.New("invalid patch file")
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(outBytes)
	if err != nil {
		return err
	}

	return nil
}

func ApplyPatchFile(patchPath string, inPath string, outPath string) error {
	compressed, err := os.ReadFile(patchPath)
	if err != nil {
		return err
	}

	uncompressed, err := decompress(compressed)
	if err != nil {
		return err
	}

	patch := Patch{}

	err = Decode[Patch](uncompressed, &patch)
	if err != nil {
		return err
	}

	inBytes, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	outBytes := []byte{}

	if xxhash.Sum64(inBytes) != patch.FromHash {
		return errors.New("invalid input file")
	}

	for _, chunk := range patch.Chunks {
		if chunk.IsBlock {
			outBytes = append(outBytes, inBytes[chunk.Block.Start:chunk.Block.End]...)
		} else {
			outBytes = append(outBytes, chunk.Data...)
		}
	}

	if xxhash.Sum64(outBytes) != patch.ToHash {
		return errors.New("invalid patch file")
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(outBytes)
	if err != nil {
		return err
	}

	return nil
}

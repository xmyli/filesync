package delta

import (
	"math"
	"os"

	"github.com/cespare/xxhash/v2"
)

func GetSignature(path string) (Signature, error) {
	signature := Signature{}

	data, err := os.ReadFile(path)
	if err != nil {
		return signature, err
	}

	signature.Hash = xxhash.Sum64(data)
	signature.BlockSize = getBlockSize(len(data))
	signature.BlockHashes = map[uint64]BlockInfo{}

	for start := 0; start < len(data); start += signature.BlockSize {
		end := start + signature.BlockSize
		if end > len(data) {
			end = len(data)
		}
		block := data[start:end]
		hash := xxhash.Sum64(block)
		signature.BlockHashes[hash] = BlockInfo{Start: start, End: end}
	}

	return signature, nil
}

func getBlockSize(sizeInBytes int) int {
	// return 8
	if sizeInBytes <= 490000 {
		return 700
	}
	blockSize := math.Sqrt(float64(sizeInBytes))
	blockSize = math.Round(blockSize/8) * 8
	if blockSize > 131072 {
		return 131072
	}
	return int(blockSize)
}

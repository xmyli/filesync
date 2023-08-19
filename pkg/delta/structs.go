package delta

type Signature struct {
	Hash        uint64
	BlockSize   int
	BlockHashes map[uint64]BlockInfo
}

type BlockInfo struct {
	Start int
	End   int
}

type Patch struct {
	FromHash uint64
	ToHash   uint64
	Chunks   []Chunk
}

type Chunk struct {
	IsBlock bool
	Block   BlockInfo
	Data    []byte
}

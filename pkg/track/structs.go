package track

type FSData struct {
	Objects map[string]FSObject
}

type FSObject struct {
	IsDir     bool
	Size      int64
	ModTime   int64
	Hash      uint64
	Signature []byte
}

func (a FSObject) SameAs(b FSObject) bool {
	return a.IsDir == b.IsDir && a.Size == b.Size && a.ModTime == b.ModTime && a.Hash == b.Hash
}

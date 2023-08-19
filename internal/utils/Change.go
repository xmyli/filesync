package utils

type ChangeType uint8

const (
	Create ChangeType = iota
	Delete
	Update
	Move
	Copy
)

type Change struct {
	Type          ChangeType
	IsDir         bool
	FromPath      string
	ToPath        string
	FromHash      uint64
	ToHash        uint64
	FromSignature []byte
	ToSignature   []byte
	Patch         []byte
	ModTime       int64
}

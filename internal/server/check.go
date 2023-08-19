package server

import (
	"filesync/internal/utils"
	"filesync/pkg/delta"
	"log"
	"os"

	"github.com/cespare/xxhash/v2"
)

func CheckBefore(root string, change utils.Change) bool {
	switch change.Type {
	case utils.Create:
		return beforeCreate(root, change)
	case utils.Delete:
		return beforeDelete(root, change)
	case utils.Update:
		return beforeUpdate(root, change)
	case utils.Move:
		return beforeMove(root, change)
	case utils.Copy:
		return beforeCopy(root, change)
	default:
		return false
	}
}

func CheckAfter(root string, change utils.Change) bool {
	switch change.Type {
	case utils.Create:
		return afterCreate(root, change)
	case utils.Delete:
		return afterDelete(root, change)
	case utils.Update:
		return afterUpdate(root, change)
	case utils.Move:
		return afterMove(root, change)
	case utils.Copy:
		return afterCopy(root, change)
	default:
		return false
	}
}

func beforeCreate(root string, change utils.Change) bool {
	if doesNotExist(root, change.ToPath) {
		return true
	}
	return false
}

func beforeDelete(root string, change utils.Change) bool {
	if exists(root, change.FromPath) && sameSignature(root, change.IsDir, change.FromPath, change.FromSignature) {
		return true
	}
	return false
}

func beforeUpdate(root string, change utils.Change) bool {
	if exists(root, change.FromPath) && sameSignature(root, change.IsDir, change.FromPath, change.FromSignature) {
		return true
	}
	return false
}

func beforeMove(root string, change utils.Change) bool {
	if exists(root, change.FromPath) && sameSignature(root, change.IsDir, change.FromPath, change.FromSignature) && doesNotExist(root, change.ToPath) {
		return true
	}
	return false
}

func beforeCopy(root string, change utils.Change) bool {
	if exists(root, change.FromPath) && sameSignature(root, change.IsDir, change.FromPath, change.FromSignature) && doesNotExist(root, change.ToPath) {
		return true
	}
	return false
}

func afterCreate(root string, change utils.Change) bool {
	if exists(root, change.ToPath) && sameSignature(root, change.IsDir, change.ToPath, change.ToSignature) {
		return true
	}
	return false
}

func afterDelete(root string, change utils.Change) bool {
	if doesNotExist(root, change.FromPath) {
		return true
	}
	return false
}

func afterUpdate(root string, change utils.Change) bool {
	if exists(root, change.ToPath) && sameSignature(root, change.IsDir, change.ToPath, change.ToSignature) {
		return true
	}
	return false
}

func afterMove(root string, change utils.Change) bool {
	if exists(root, change.ToPath) && sameSignature(root, change.IsDir, change.ToPath, change.ToSignature) {
		return true
	}
	return false
}

func afterCopy(root string, change utils.Change) bool {
	if exists(root, change.FromPath) && exists(root, change.ToPath) && sameSignature(root, change.IsDir, change.FromPath, change.FromSignature) && sameSignature(root, change.IsDir, change.ToPath, change.ToSignature) {
		return true
	}
	return false
}

func exists(root string, path string) bool {
	if _, err := os.Stat(root + "/" + path); err != nil {
		return false
	}
	return true
}

func doesNotExist(root string, path string) bool {
	_, err := os.Stat(root + "/" + path)
	if err == nil {
		return false
	}
	if !os.IsNotExist(err) {
		return false
	}
	return true
}

func sameSignature(root string, isDir bool, aPath string, bSignatureBytes []byte) bool {
	if isDir {
		return isDir
	}

	aSignature, err := delta.GetSignature(root + "/" + aPath)
	if err != nil {
		log.Fatalln(err)
	}

	aSignatureBytes, err := delta.Encode(aSignature)
	if err != nil {
		log.Fatalln(err)
	}

	aSignatureHash := xxhash.Sum64(aSignatureBytes)
	bSignatureHash := xxhash.Sum64(bSignatureBytes)

	return aSignatureHash == bSignatureHash
}

package utils

import (
	"filesync/pkg/delta"
	"filesync/pkg/track"
	"log"

	"github.com/cespare/xxhash/v2"
)

func Compare(before map[string]track.FSObject, after map[string]track.FSObject) []Change {
	changes := []Change{}

	added := []string{}
	removed := []string{}

	// Find added
	for path, data := range after {
		_, exists := before[path]
		if !exists {
			change := Change{
				Type:          Create,
				IsDir:         data.IsDir,
				FromPath:      "",
				ToPath:        path,
				FromHash:      0,
				ToHash:        data.Hash,
				FromSignature: []byte{},
				ToSignature:   data.Signature,
				ModTime:       data.ModTime,
			}
			changes = append(changes, change)
			added = append(added, path)
		}
	}

	// Find removed
	for path, data := range before {
		_, exists := after[path]
		if !exists {
			change := Change{
				Type:          Delete,
				IsDir:         data.IsDir,
				FromPath:      path,
				ToPath:        "",
				FromHash:      data.Hash,
				ToHash:        0,
				FromSignature: data.Signature,
				ToSignature:   []byte{},
				ModTime:       data.ModTime,
			}
			changes = append(changes, change)
			removed = append(removed, path)
		}
	}

	// Find edited
	for path, beforeData := range before {
		afterData, exists := after[path]
		if exists && !afterData.IsDir && !beforeData.SameAs(afterData) {
			change := Change{
				Type:          Update,
				IsDir:         afterData.IsDir,
				FromPath:      path,
				ToPath:        path,
				FromHash:      beforeData.Hash,
				ToHash:        afterData.Hash,
				FromSignature: beforeData.Signature,
				ToSignature:   afterData.Signature,
				ModTime:       max(beforeData.ModTime, afterData.ModTime),
			}
			changes = append(changes, change)
		}
	}

	// Find moved
	reversedAddedFiles := map[uint64]string{}
	for _, path := range added {
		fileData := after[path]
		if !fileData.IsDir {
			hashBytes, err := delta.Encode(fileData)
			if err != nil {
				log.Fatalln(err)
			}
			hash := xxhash.Sum64(hashBytes)
			reversedAddedFiles[hash] = path
		}
	}
	for _, oldPath := range removed {
		fileData := before[oldPath]
		hashBytes, err := delta.Encode(fileData)
		if err != nil {
			log.Fatalln(err)
		}
		hash := xxhash.Sum64(hashBytes)
		newPath, exists := reversedAddedFiles[hash]
		if exists {
			change := Change{
				Type:          Move,
				IsDir:         fileData.IsDir,
				FromPath:      oldPath,
				ToPath:        newPath,
				FromHash:      fileData.Hash,
				ToHash:        fileData.Hash,
				FromSignature: fileData.Signature,
				ToSignature:   fileData.Signature,
				ModTime:       fileData.ModTime,
			}
			changes = append(changes, change)
			removed = remove(removed, oldPath)
			added = remove(added, newPath)
		}
	}

	// Find copied
	reversedUnchangedFiles := map[uint64]string{}
	for path, beforeData := range before {
		afterData, samePath := after[path]

		if samePath && !afterData.IsDir && beforeData.SameAs(afterData) {
			hashBytes, err := delta.Encode(afterData)
			if err != nil {
				log.Fatalln(err)
			}
			hash := xxhash.Sum64(hashBytes)
			reversedUnchangedFiles[hash] = path
		}
	}
	for _, copyPath := range added {
		fileData := after[copyPath]
		hashBytes, err := delta.Encode(fileData)
		if err != nil {
			log.Fatalln(err)
		}
		hash := xxhash.Sum64(hashBytes)
		originalPath, exists := reversedUnchangedFiles[hash]
		if exists {
			change := Change{
				Type:          Copy,
				IsDir:         fileData.IsDir,
				FromPath:      originalPath,
				ToPath:        copyPath,
				FromHash:      fileData.Hash,
				ToHash:        fileData.Hash,
				FromSignature: fileData.Signature,
				ToSignature:   fileData.Signature,
				ModTime:       fileData.ModTime,
			}
			changes = append(changes, change)
			added = remove(added, copyPath)
		}
	}

	return changes
}

func remove(s []string, e string) []string {
	for i, v := range s {
		if e == v {
			s[i] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

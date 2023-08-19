package utils

import (
	"filesync/pkg/delta"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Process(root string, changes []Change) {
	err := os.MkdirAll(filepath.Base(root)+"/.trash", 0777)
	if err != nil {
		log.Fatalln(err)
	}

	moves := map[string]string{}

	for _, change := range changes {
		switch change.Type {
		case Create:
			tempFilePath, toPath := processCreate(root, change)
			moves[tempFilePath] = filepath.Clean(toPath)
		case Delete:
			moves[change.FromPath] = filepath.Base(root) + "/.trash"
		case Update:
			tempFilePath, toPath := processUpdate(root, change)
			moves[tempFilePath] = filepath.Clean(toPath)
		case Move:
			moves[change.FromPath] = filepath.Clean(change.ToPath)
		case Copy:
			tempFilePath, toPath := processCopy(root, change)
			moves[tempFilePath] = filepath.Clean(toPath)
		}
	}

	for from, to := range moves {
		err := os.Rename(from, to)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func processCreate(root string, change Change) (string, string) {
	toPath := root + "/" + change.ToPath

	tempFilePath := os.TempDir() + strconv.Itoa(int(rand.Int63()))

	if change.IsDir {
		err := os.Mkdir(tempFilePath, 0777)
		if err != nil {
			log.Fatalln(err)
		}
		return tempFilePath, toPath
	}

	empty, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.Remove(empty.Name())

	err = delta.ApplyPatch(change.Patch, empty.Name(), tempFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	modTime := time.UnixMicro(change.ModTime)

	err = os.Chtimes(tempFilePath, modTime, modTime)
	if err != nil {
		log.Fatalln(err)
	}

	return tempFilePath, toPath
}

func processUpdate(root string, change Change) (string, string) {
	fromPath := root + "/" + change.FromPath
	toPath := root + "/" + change.ToPath

	tempFilePath := os.TempDir() + strconv.Itoa(int(rand.Int63()))

	err := delta.ApplyPatch(change.Patch, fromPath, tempFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	modTime := time.UnixMicro(change.ModTime)

	err = os.Chtimes(tempFilePath, modTime, modTime)
	if err != nil {
		log.Fatalln(err)
	}

	return tempFilePath, toPath
}

func processCopy(root string, change Change) (string, string) {
	fromPath := root + "/" + change.FromPath
	toPath := root + "/" + change.ToPath

	tempFilePath := os.TempDir() + strconv.Itoa(int(rand.Int63()))

	file, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := os.ReadFile(fromPath)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalln(err)
	}

	modTime := time.UnixMicro(change.ModTime)

	err = os.Chtimes(tempFilePath, modTime, modTime)
	if err != nil {
		log.Fatalln(err)
	}

	return tempFilePath, toPath
}

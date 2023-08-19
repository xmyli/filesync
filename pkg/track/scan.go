package track

import (
	"filesync/pkg/delta"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func (t *Tracker) walk(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		log.Fatal(err)
	}

	if path == t.Path {
		return nil
	}

	if slices.Contains(t.ignoredNames, filepath.Base(path)) {
		return nil
	}

	info, err := entry.Info()
	if err != nil {
		log.Fatal(err)
	}

	key := strings.TrimPrefix(path, filepath.Clean(t.Path))

	object, exists := t.currentState.Objects[key]
	if exists {
		same := object.IsDir == info.IsDir() && object.Size == info.Size() && object.ModTime == info.ModTime().UnixMicro()
		if same {
			t.nextState.Objects[key] = object
			return nil
		}
	}

	if info.IsDir() {
		t.nextState.Objects[key] = FSObject{
			IsDir:     info.IsDir(),
			Size:      info.Size(),
			ModTime:   info.ModTime().UnixMicro(),
			Hash:      0,
			Signature: []byte{},
		}
		return nil
	}

	signature, err := delta.GetSignature(path)
	if err != nil {
		log.Fatal(err)
	}

	signatureBytes, err := delta.Encode(signature)
	if err != nil {
		log.Fatal(err)
	}

	t.nextState.Objects[key] = FSObject{
		IsDir:     info.IsDir(),
		Size:      info.Size(),
		ModTime:   info.ModTime().UnixMicro(),
		Hash:      signature.Hash,
		Signature: signatureBytes,
	}

	return nil
}

func (t *Tracker) scan() {
	filepath.WalkDir(t.Path, t.walk)
	t.currentState = t.nextState
	t.nextState = FSData{
		Objects: map[string]FSObject{},
	}
	log.Println("SCAN")
	t.updateTimer.Reset(t.updateInterval)
	t.onScan(t.currentState)
}

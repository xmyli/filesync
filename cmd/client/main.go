package main

import (
	"filesync/internal/client"
	"filesync/internal/utils"
	"filesync/pkg/track"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dirPtr := flag.String("d", "", "dir")

	flag.Parse()

	if *dirPtr == "" {
		log.Fatalln("dir not specified")
	}

	rpcClient := client.RPCClient{}
	rpcClient.Root = filepath.Clean(*dirPtr)

	if info, err := os.Stat(rpcClient.Root); !os.IsNotExist(err) {
		if info.IsDir() {
			lastSyncId := ""
			lastScanState := map[string]track.FSObject{}
			lastSyncState := map[string]track.FSObject{}
			initialized := false

			// File System Tracker
			tracker := track.NewTracker(rpcClient.Root, time.Minute, time.Second, []string{".trash"})
			tracker.SetOnScan(func(f track.FSData) {
				if !initialized {
					lastSyncState, lastSyncId = client.Sync(&rpcClient, f, f.Objects, f.Objects, []utils.Change{}, "") // should load from disk or if empty, clone server state
					lastScanState = lastSyncState
					initialized = true
				}

				clientChanges := utils.Compare(lastScanState, f.Objects)
				lastScanState = f.Objects
				if len(clientChanges) > 0 {
					lastSyncState, lastSyncId = client.Sync(&rpcClient, f, lastSyncState, lastScanState, clientChanges, lastSyncId)
					lastScanState = lastSyncState
				}
			})
			go tracker.Start()

			// SSE Client
			sseDataCh := make(chan []byte)
			sseClient := client.NewSSEClient(sseDataCh)
			go sseClient.Connect("http://127.0.0.1:8123")
			for {
				data := <-sseDataCh
				if string(data) != lastSyncId {
					lastSyncState, lastSyncId = client.Sync(&rpcClient, tracker.GetState(), lastSyncState, lastScanState, []utils.Change{}, lastSyncId)
					lastScanState = lastSyncState
				}
			}
		} else {
			log.Fatalln("Not a folder.")
		}
	} else {
		log.Fatalln("Path does not exist.")
	}
}

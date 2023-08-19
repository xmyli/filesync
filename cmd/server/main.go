package main

import (
	"filesync/internal/server"
	"filesync/internal/utils"
	"filesync/pkg/track"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dirPtr := flag.String("d", "", "dir")

	flag.Parse()

	if *dirPtr == "" {
		log.Fatalln("dir not specified")
	}

	if info, err := os.Stat(*dirPtr); !os.IsNotExist(err) {
		if info.IsDir() {
			lastScanState := map[string]track.FSObject{}

			// SSE Server
			sseServer := server.NewSSEServer()
			go http.ListenAndServe(":8123", sseServer)

			// File System Tracker
			tracker := track.NewTracker(*dirPtr, time.Minute, time.Second, []string{".trash"})
			tracker.SetOnScan(func(f track.FSData) {
				changes := utils.Compare(lastScanState, f.Objects)
				lastScanState = f.Objects
				if len(changes) > 0 {
					sseServer.SendData([]byte{1})
				}
			})
			go tracker.Start()

			// RPC Server
			rpcServer := server.NewRPCServer(&tracker)
			go rpcServer.Start()

			select {}
		} else {
			log.Fatalln("Not a folder.")
		}
	} else {
		log.Fatalln("Path does not exist.")
	}
}

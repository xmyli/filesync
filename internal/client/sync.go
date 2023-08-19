package client

import (
	"filesync/internal/utils"
	"filesync/pkg/delta"
	"filesync/pkg/track"
	"log"
)

func Sync(rpcClient *RPCClient, currentState track.FSData, lastSyncState map[string]track.FSObject, lastScanState map[string]track.FSObject, clientChanges []utils.Change, lastSyncId string) (map[string]track.FSObject, string) {
	serverState, err := rpcClient.GetState("127.0.0.1:8321")
	if err != nil {
		log.Println(err)
		return lastSyncState, lastSyncId
	}

	if len(clientChanges) == 0 {
		clientChanges = utils.Compare(lastScanState, currentState.Objects)
	}

	serverChanges := utils.Compare(lastSyncState, serverState.Objects)

	if len(serverChanges) > 0 || len(clientChanges) > 0 {
		for i := 0; i < len(clientChanges); i++ {
			change := clientChanges[i]

			if change.IsDir {
				continue
			}

			if change.Type == utils.Create {
				signature := delta.Signature{}
				signature.BlockHashes = map[uint64]delta.BlockInfo{}

				patch, err := delta.CreatePatch(signature, rpcClient.Root+change.ToPath)
				if err != nil {
					log.Fatalln(err)
				}

				clientChanges[i].Patch = patch
			} else if change.Type == utils.Update {
				signature := delta.Signature{}

				err := delta.Decode[delta.Signature](change.FromSignature, &signature)
				if err != nil {
					log.Fatalln(err)
				}

				patch, err := delta.CreatePatch(signature, rpcClient.Root+change.ToPath)
				if err != nil {
					log.Fatalln(err)
				}

				clientChanges[i].Patch = patch
			}
		}

		currentSyncId, err := rpcClient.DoUpdate("127.0.0.1:8321", clientChanges, serverChanges)
		if err != nil {
			log.Println(err)
			return lastSyncState, lastSyncId
		}

		return currentState.Objects, currentSyncId
	}

	return lastSyncState, lastSyncId
}

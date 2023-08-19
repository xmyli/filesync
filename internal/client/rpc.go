package client

import (
	"filesync/internal/server"
	"filesync/internal/utils"
	"filesync/pkg/track"
	"net/rpc"
	"strconv"
)

type RPCClient struct {
	Root string
}

func (r *RPCClient) GetState(addr string) (track.FSData, error) {
	args := server.GetStateArgs{}
	reply := server.GetStateReply{}

	err := call(addr, "RPCServer.GetState", &args, &reply)
	if err != nil {
		return track.FSData{}, err
	}

	return reply.State, nil
}

func (r *RPCClient) DoUpdate(addr string, clientChanges []utils.Change, serverChanges []utils.Change) (string, error) {
	args := server.DoUpdateArgs{}

	args.ClientChanges = clientChanges
	args.ServerChanges = serverChanges

	reply := server.DoUpdateReply{}

	err := call(addr, "RPCServer.DoUpdate", &args, &reply)
	if err != nil {
		return "", err
	}

	utils.Process(r.Root, reply.ServerChanges)

	return strconv.Itoa(int(reply.Id)), nil
}

func call(addr string, service string, args interface{}, reply interface{}) error {
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}

	err = rpcClient.Call(service, args, reply)
	if err != nil {
		rpcClient.Close()
		return err
	}

	rpcClient.Close()
	return nil
}

package server

import (
	"filesync/internal/utils"
	"math/rand"
)

type DoUpdateArgs struct {
	ClientChanges []utils.Change
	ServerChanges []utils.Change
}

type DoUpdateReply struct {
	Id            int64
	Rejected      bool
	ServerChanges []utils.Change
}

func (s *RPCServer) DoUpdate(args *DoUpdateArgs, reply *DoUpdateReply) error {
	for _, change := range args.ClientChanges {
		if !CheckBefore(s.tracker.Path, change) {
			reply.Rejected = true
			return nil
		}
	}

	for _, change := range args.ServerChanges {
		if !CheckAfter(s.tracker.Path, change) {
			reply.Rejected = true
			return nil
		}
	}

	utils.Process(s.tracker.Path, args.ClientChanges)

	toClient := []utils.Change{}
	for _, change := range args.ServerChanges {
		toClient = append(toClient, populate(s.tracker.Path, change))
	}

	reply.Id = rand.Int63()
	reply.Rejected = false
	reply.ServerChanges = toClient

	return nil
}

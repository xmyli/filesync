package server

import (
	"filesync/pkg/track"
)

type GetStateArgs struct {
}

type GetStateReply struct {
	State track.FSData
}

func (s *RPCServer) GetState(args *GetStateArgs, reply *GetStateReply) error {
	reply.State = s.tracker.GetState()

	return nil
}

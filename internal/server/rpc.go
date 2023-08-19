package server

import (
	"filesync/pkg/track"
	"log"
	"net"
	"net/rpc"
)

type RPCServer struct {
	tracker *track.Tracker
}

func (s *RPCServer) Start() {
	rpcs := rpc.NewServer()
	rpcs.Register(s)
	listener, err := net.Listen("tcp", ":8321")
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			} else {
				break
			}
		}
		listener.Close()
	}()
}

func NewRPCServer(tracker *track.Tracker) *RPCServer {
	return &RPCServer{tracker: tracker}
}

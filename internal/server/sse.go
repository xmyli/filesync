package server

import (
	"context"
	"log"
	"net/http"
	"sync"
)

type Connection struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	ctx     context.Context
}

type SSEServer struct {
	mu          sync.RWMutex
	connections map[string]*Connection
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		connections: map[string]*Connection{},
	}
}

func (s *SSEServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		log.Fatalln("Flusher not supported")
	}

	requestContext := req.Context()
	s.mu.Lock()
	s.connections[req.RemoteAddr] = &Connection{
		writer:  rw,
		flusher: flusher,
		ctx:     requestContext,
	}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.connections, req.RemoteAddr)
		s.mu.Unlock()
	}()

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	<-requestContext.Done()
}

func (s *SSEServer) SendData(data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client, connection := range s.connections {
		_, err := connection.writer.Write(data)
		if err != nil {
			s.mu.Lock()
			delete(s.connections, client)
			s.mu.Unlock()
			continue
		}

		connection.flusher.Flush()
	}
}

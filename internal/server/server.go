package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

const protocol = "tcp"

func Serve(port int, handler Handler) (*Server, error) {

	listener, err := net.Listen(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error creating listener via %s on %d", protocol, port)
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}
	// server = running
	server.closed.Store(false)

	go server.listen()

	return server, nil

}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("Error closing listener")
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting new connection")
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	writer := response.NewWriter(conn)

	request, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.BadRequest)
		body := []byte(err.Error())
		writer.WriteHeaders(response.GetDefaultHeaders(len(body)))
		writer.WriteBody(body)
		return
	}
	s.handler(writer, request)
}

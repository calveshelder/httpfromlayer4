package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/calveshelder/httpfromlayer4/internal/response"
)

type Server struct {
	state     int
	listening atomic.Bool
	listener  net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: l}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.listening.Store(false)
	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.listening.Load() {
				return
			}
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	_ = response.WriteStatusLine(conn, response.StatusOk)
	headers := response.GetDefaultHeaders(0)
	_ = response.WriteHeaders(conn, headers)
}

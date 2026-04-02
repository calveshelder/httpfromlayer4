package server

import (
	"fmt"
	"net"
	"sync/atomic"
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
	status := "HTTP/1.1 200 OK"
	contentType := "text/plain"
	contentLength := "13"
	contentBody := "Hello World!"

	conn.Write([]byte(fmt.Sprintf("%s\r\n%sContent-Type: \r\n%sContent-Length: \r\n\r\n%s", status, contentType, contentLength, contentBody)))
	conn.Close()
}

package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/calveshelder/httpfromlayer4/internal/request"
	"github.com/calveshelder/httpfromlayer4/internal/response"
)

type Server struct {
	state     int
	listening atomic.Bool
	listener  net.Listener
	handler   Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Status  int
	Message string
}

func Serve(port int, f Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: l, handler: f}

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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println(err)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}
	_ = response.WriteStatusLine(conn, response.StatusOk)
	headers := response.GetDefaultHeaders(buf.Len())
	_ = response.WriteHeaders(conn, headers)
	_, _ = conn.Write(buf.Bytes())
}

func (he HandlerError) Write(w io.Writer) error {
	if err := response.WriteStatusLine(w, response.StatusCode(he.Status)); err != nil {
		return err
	}
	headers := response.GetDefaultHeaders(len(he.Message))
	if err := response.WriteHeaders(w, headers); err != nil {
		return err
	}
	if _, err := w.Write([]byte(he.Message)); err != nil {
		return err
	}

	return nil
}

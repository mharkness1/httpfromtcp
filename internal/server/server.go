package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/mharkness1/httpfromtcp/internal/request"
	"github.com/mharkness1/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  HandlerFunc
}

func Serve(port int, handler HandlerFunc) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
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
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		handlerErr.HandlerErrorWriter(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buf, req)
	if handlerError != nil {
		handlerError.HandlerErrorWriter(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}

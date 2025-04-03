package server

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type State int

type Server struct {
	currentState State
	port         int
	listener     net.Listener
	conn         io.Writer
	handler			 Handler
}

const (
	ServerListeningState State = iota
	ServerClosedState
)

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s := Server{
		currentState: ServerListeningState,
		port:         port,
		listener:     l,
		handler: handler,
	}

	go s.listen()
	return &s, nil
}

func (s *Server) Close() error {
	if s.currentState == ServerClosedState {
		return errors.New("server is already closed")
	}
	s.currentState = ServerClosedState
	return s.listener.Close()
}

func (s *Server) listen() {
	for s.currentState != ServerClosedState {
		conn, err := s.listener.Accept()
		s.conn = conn
		if err != nil {
			if s.currentState != ServerClosedState {
				log.Println(err)
			}
			continue
		}
		go s.handle(conn)
	}
}

func WriteHandlerError(w io.Writer, err *HandlerError) {
	if err != nil {
		fmt.Fprintf(w, "Error (%d): %s\n", err.StatusCode, err.Message)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}

	b := new(bytes.Buffer)

	hErr := s.handler(b, r)
	if hErr != nil {
		response.WriteStatusLine(conn, response.StatusCode(hErr.StatusCode))
		conn.Write([]byte("\r\n" + hErr.Message))
		return
	}

	response.WriteStatusLine(conn, response.StatusCode(200))
	headers := response.GetDefaultHeaders(len(b.Bytes()))
	response.WriteHeaders(conn, headers)

	conn.Write(b.Bytes())

	// if err != nil {
	// 	fmt.Println(err)
	// 	s.currentState = ServerClosedState
	// 	return
	// }
}

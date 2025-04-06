package server

import (
	"bytes"
	"errors"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
)

type Handler func(w response.Writer, req *request.Request)

type State int

type Server struct {
	currentState State
	port         int
	listener     net.Listener
	Conn         io.Writer
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
		s.Conn = conn
		if err != nil {
			if s.currentState != ServerClosedState {
				log.Println(err)
			}
			continue
		}
		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// r, err := request.RequestFromReader(conn)
	// if err != nil {
	// 	return
	// }

	b := new(bytes.Buffer)

	// response.WriteStatusLine(response.StatusCode(200))
	// headers := response.GetDefaultHeaders(len(b.Bytes()))
	// response.WriteHeaders(headers)
	// response.WriteB

	conn.Write(b.Bytes())
}

package server

import (
	"errors"
	"fmt"
	// "io"
	"log"
	"net"
	"strconv"
)

type State int

type Server struct{
	currentState State
	port int
	listener net.Listener
}

const (
	ServerListeningState State = iota
	ServerClosedState
)

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s := Server{
		currentState: ServerListeningState,
		port: port,
		listener: l,
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
			if err != nil {
					// Only log if we're not closed - errors are expected when closing
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
	resp := "HTTP/1.1 200 OK\r\n" +
                "Content-Type: text/plain\r\n" +
                "\r\n" +
                "Hello World!"

	_, err := conn.Write([]byte(resp))
	if err != nil {
		fmt.Println(err)
		s.currentState = ServerClosedState
		return
	}
}
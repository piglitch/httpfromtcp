package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func HandlerFunc(w response.Writer, req *request.Request) {
	w.Write = new(io.Writer)
	if req.RequestLine.RequestTarget == "/yourproblem"{
		
	}
	if req.RequestLine.RequestTarget == "/myproblem"{
		
	}
}

func main() {
	server, err := server.Serve(port, HandlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
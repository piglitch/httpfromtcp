package main

import (
	"fmt"
	"log"
	"net"
	"httpfromtcp/internal/request"
)

func main(){
	address := ":42069"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not listen to port number: %s. Error: %s", address, err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			return
		}
		fmt.Println("A connection has been accepted.")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error: ", err)
			return
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", r.RequestLine.Method)
		fmt.Println("- Target:", r.RequestLine.RequestTarget)
		fmt.Println("- Version:", r.RequestLine.HttpVersion)
	}
}
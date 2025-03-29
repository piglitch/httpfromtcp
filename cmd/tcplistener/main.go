package main

import (
	"fmt"
	"log"
	"net"
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
		msgChan := getLinesChannel(conn)
		for msg := range msgChan {
			fmt.Println(msg)
		}		
		fmt.Println("Channel is closed")
	}
}
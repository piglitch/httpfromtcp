package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	// "os"
)

func main(){
	address := ":42069"
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Could not listen to port number: %s. Error: %s", address, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		msg, _ := r.ReadString('\n')
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Println("errors: ", err)
			return
		}
	}
}
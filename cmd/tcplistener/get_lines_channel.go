package main

import (
	// "fmt"
	// "errors"
	// "io"
	"net"

	// "log"
	// "os"
	"strings"
)

func getLinesChannel(conn net.Conn) <-chan string{
	msg := make(chan string)
	b := make([]byte, 8)
	var currLine string 
	
	go func() {
		defer close(msg)
		for {
			byte_num, err := conn.Read(b)
			if err != nil {
				msg <- currLine
				return
			}
			currLine += string(b[:byte_num])
			parts := strings.Split(currLine, "\n")

			for i, part := range parts {
				if i < len(parts)-1 {
					msg <- part	
				}
			}
			currLine = parts[len(parts)-1]
		}	
	}()
	return msg
}
package main

import (
	// "fmt"
	"errors"
	"io"
	// "log"
	// "os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string{
	msg := make(chan string)
	b := make([]byte, 8)
	var currLine string 
	
	go func() {
		defer close(msg)
		defer f.Close()
		for {
			byte_num, err := f.Read(b)
			if errors.Is(err, io.EOF){
				return
			}
			currLine += string(b[:byte_num])
			parts := strings.Split(currLine, "\n")

			for i, part := range parts {
				if i < len(parts)-1 && part != "" {
					msg <- part	
				}
			}
			currLine = parts[len(parts)-1]
		}	
	}()

	return msg
}
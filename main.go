package main

import (
	"fmt"
	"log"
	"os"
)

func main(){
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("could not open the file", err)
		return 
	}
	msgChan := getLinesChannel(file)

	for msg := range msgChan {
		fmt.Printf("read: %s\n", string(msg))
	}
}
package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// b := make([]byte, 8)
	b, err := io.ReadAll(reader)
	// n, err := reader.Read(b)
	if err != nil {
		log.Fatal("Failed to read")
		return nil, err
	}
	requestString := string(b)
	reqSlice := strings.Split(requestString, "\r\n")
	reqLineString := reqSlice[0]
	reqLineSlice := strings.Split(reqLineString, " ") 

	if len(reqLineSlice) < 3 {
		return nil, errors.New("poorly formatted request")
	}

	method := reqLineSlice[0]
	if method != strings.ToUpper(method) || method == "" {	
		return nil, errors.New("method cannot be in lowercase")
	}
	version := reqLineSlice[2][5:]
	if version != "1.1" && reqLineSlice[2] != "HTTP/1.1" {
		fmt.Println(version, "line 39 version")
		newError := "only HTTP/1.1 is allowed: " + version + " " + method
		return nil, errors.New(newError)
	}

	requestLine := RequestLine{
		HttpVersion: version,
		RequestTarget: reqLineSlice[1],
		Method: method,
	}

	r := Request{
		RequestLine: requestLine,
	}

	return &r, nil 
}
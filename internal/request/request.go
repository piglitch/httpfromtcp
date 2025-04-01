package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"httpfromtcp/internal/headers"
)

type parseState int

const (
	initialized parseState = iota
	done 
	requestStateParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	parseState int
	Headers map[string]string
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

var h headers.Headers = make(headers.Headers)

var ErrNeedMoreData = errors.New("need more data")

var CRLF string = "\r\n"

func parseRequestLine(b []byte) (r RequestLine, n int, err error) {
	requestString := string(b)
	crlf_idx := strings.Index(requestString, CRLF)
	requestString = requestString[:crlf_idx+2]
	reqSlice := strings.Split(requestString, "\r\n")
	if len(reqSlice) < 2 {
		return RequestLine{}, 0, nil
	}
	reqLineString := reqSlice[0]
	reqLineSlice := strings.Split(reqLineString, " ") 

	if len(reqLineSlice) < 3 {
		return RequestLine{}, 0, errors.New("poorly formatted request" + reqLineString)
	}

	method := reqLineSlice[0]
	if method != strings.ToUpper(method) || method == "" {	
		return RequestLine{}, 0, errors.New("method cannot be in lowercase")
	}
	version := reqLineSlice[2][5:]
	if version != "1.1" && reqLineSlice[2] != "HTTP/1.1" {
		fmt.Println(version, "line 39 version")
		newError := "only HTTP/1.1 is allowed: " + version + " " + method
		return RequestLine{}, 0, errors.New(newError)
	}

	requestLine := RequestLine{
		HttpVersion: version,
		RequestTarget: reqLineSlice[1],
		Method: method,
	}

	return requestLine, len(reqLineString) + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	if r.parseState == int(initialized) {
		rL, n, err := parseRequestLine(data)	
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, ErrNeedMoreData
		}
		r.RequestLine = rL
		totalBytesParsed += n
		r.parseState = int(requestStateParsingHeaders)
		return n, nil
	}
	if r.parseState == int(done) {
		return 0, errors.New("error: trying to read data in a done state")
	}
	for r.parseState == int(requestStateParsingHeaders) {
		n, state, err := h.Parse(data[totalBytesParsed:])
		if err == nil {
			return 0, err
		}
		totalBytesParsed += n
		if state || totalBytesParsed >= len(data[totalBytesParsed:]) {
			r.parseState = int(done)
			return n, nil
		}
		continue
	}
	return 0, errors.New("error: unknown state")
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	bufferSize := 8
	b := make([]byte, bufferSize)
	readToIndex := 0

	r := &Request{
		parseState: int(initialized),
	}

	for r.parseState != int(done) {

		if readToIndex == len(b) {
			new_b := make([]byte, len(b) * 2)
			copy(new_b, b)
			b = new_b
		}
		
		rn, err := reader.Read(b[readToIndex:])

		if err == io.EOF {
			r.parseState = int(done)
			break
		}
		readToIndex += rn
		pn, err := r.parse(b[:readToIndex])
		if errors.Is(err, ErrNeedMoreData) {
			continue
		}
		if !errors.Is(err, ErrNeedMoreData) && err != nil {
			return &Request{}, err
		}
	
		copy(b, b[pn:])
		readToIndex -= pn
		
	}

	return r, nil 
}
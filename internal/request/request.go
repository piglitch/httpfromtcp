package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	state requestState
	Headers headers.Headers
	Body []byte
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

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

func (r *Request) Parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.state == requestStateInitialized {
		rL, n, err := parseRequestLine(data)	
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, ErrNeedMoreData
		}
		r.RequestLine = rL
		r.state = requestStateParsingHeaders
		return n, nil
	}

	if r.state == requestStateParsingHeaders {
		n, state, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if state {
			r.state = requestStateParsingBody
		}
		return n, nil
	}

	if r.state == requestStateParsingBody {
		contentLength := r.Headers.Get("Content-Length")
		if contentLength == "" {
			r.state = requestStateDone
			return len(data), nil
		}
		r.Body = append(r.Body, data...)
		contentLengthInt, err := strconv.Atoi(contentLength)
    if err != nil {
        return 0, fmt.Errorf("invalid Content-Length header: %v", err)
    }
		
		if len(r.Body) > contentLengthInt {
			return 0, errors.New("body is larger than than reported content length")
		}

		if contentLengthInt == len(r.Body) {
			r.state = requestStateDone
			return len(data), nil
		}

		return len(data), nil
	}

	if r.state == requestStateDone {
		return 0, errors.New("error: trying to read data in a done state")
	}

	return 0, errors.New("error: unknown state" + string(data))
}


func RequestFromReader(reader io.Reader) (*Request, error) {

	bufferSize := 8
	b := make([]byte, bufferSize)
	readToIndex := 0

	h := make(headers.Headers)

	r := &Request{
		state: requestStateInitialized,
		Headers: h,
	}

	for r.state != requestStateDone {

		if readToIndex == len(b) {
			new_b := make([]byte, len(b) * 2)
			copy(new_b, b)
			b = new_b
		}
		
		rn, err := reader.Read(b[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", r.state, rn)
				}
				break
			}
			return nil, err
		}
		readToIndex += rn
		pn, err := r.Parse(b[:readToIndex])
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
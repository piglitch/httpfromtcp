package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode= 400
	StatusInternalError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase := ""
	if statusCode == StatusOk {
		statusCode = StatusOk
		reasonPhrase = "HTTP/1.1 " + strconv.Itoa(int(StatusOk)) + " OK" + "\r\n"
	}
	if statusCode == 400 {
		statusCode = StatusBadRequest
		reasonPhrase = "HTTP/1.1 " + strconv.Itoa(int(StatusBadRequest)) + " Bad Request" + "\r\n"
	}
	if statusCode == 500 {
		statusCode = StatusInternalError
		reasonPhrase = "HTTP/1.1 " + strconv.Itoa(int(StatusInternalError)) + " Internal Server Error" + "\r\n"
	}

	_, err := w.Write([]byte(reasonPhrase))
	if err != nil {
		return err
	}
	return nil
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var str string
	for key := range headers {
		str += key + ": " + headers[key] + "\r\n"
	}
	str += "\r\n"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
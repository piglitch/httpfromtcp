package response

import (
	"errors"
	"fmt"

	// "net/http"
	"strconv"
	// "strings"

	// "fmt"
	"httpfromtcp/internal/headers"
	"io"
)


type StatusCode int

const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalError StatusCode = 500
)

type writerState int

const (
	stateStatus writerState = iota
	stateHeader  
	stateBody 
	stateTrailer
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)
	h["Content-Length"] = strconv.Itoa(contentLen)
	// h["Transfer-Encoding"] = "chunked"
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

type Writer struct{
	Writer 				io.Writer
	writerState 	writerState
	FullBody 			[]byte
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := ""
	if w.writerState != stateStatus {
		return errors.New("have to write status first")	
	}

	if statusCode == 200 {
		statusCode = StatusOk
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "OK")
	}

	if statusCode == 400 {
		statusCode = StatusBadRequest
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "Bad Request")
	}

	if statusCode == 500 {
		statusCode = StatusInternalError
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "Internal Server Error")
	}
	w.writerState = stateHeader
	w.Writer.Write([]byte(statusLine))

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	// if w.writerState == stateHeader {
	// 	return errors.New("can't write headers now")
	// }
	var str string
	for key := range headers {
		str += key + ": " + headers[key] + "\r\n"
	}
	str += "\r\n"

	w.Writer.Write([]byte(str))

	w.writerState = stateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	// if w.writerState == stateBody {
	// 	return 0, errors.New("can't write body now")
	// }
	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}
	w.writerState = stateTrailer
	return n, nil
}

func (w *Writer) WriterChunkedBody(p []byte) (int, error) {

	length := len(p)
	hexLength := fmt.Sprintf("%x", length)

	_, err := w.Writer.Write([]byte(hexLength + "\r\n"))
	if err != nil {
		return 0, err
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}

	_, err = w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	w.FullBody = append(w.FullBody, p...)
	return n, nil
}

func (w *Writer) WriterChunkedBodyDone() (int, error) {
	w.writerState = stateTrailer
	return w.Writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(trailers headers.Headers) error {
	var str string
	
	for key, value := range trailers {
			str += key + ": " + value + "\r\n"
	}
	str += "\r\n"
	
	fmt.Println(str, "line 140")
	
	// Write directly to the underlying writer, not to w.Writer
	_, err := w.Writer.Write([]byte(str))
	return err
}
package response

import (
	"errors"
	"fmt"
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
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)
	// h["Content-Length"] = strconv.Itoa(contentLen)
	h["Transfer-Encoding"] = "chunked"
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

type Writer struct{
	Writer 				io.Writer
	writerState writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := ""
	if w.writerState != stateStatus {
		return errors.New("have to write status first")	
	}

	if statusCode == 200 {
		statusCode = StatusOk
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "OK")
		// reasonHtml = fmt.Sprintf(`<html>
		// 														<head>
		// 															<title>%d OK</title>
		// 														</head>
		// 														<body>
		// 															<h1>Success!</h1>
		// 															<p>Your request was an absolute banger.</p>
		// 														</body>
		// 													</html>`, statusCode)
	}

	if statusCode == 400 {
		statusCode = StatusBadRequest
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "Bad Request")
		// reasonHtml = fmt.Sprintf(`<html>
		// 								<head>
		// 									<title>%d Bad Request</title>
		// 								</head>
		// 								<body>
		// 									<h1>Bad Request</h1>
		// 									<p>Your request honestly kinda sucked.</p>
		// 								</body>
		// 							</html>`, statusCode)
	}

	if statusCode == 500 {
		statusCode = StatusInternalError
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, "Internal Server Error")
		// reasonHtml = fmt.Sprintf(`<html>
		// 														<head>
		// 															<title>%d Internal Server Error</title>
		// 														</head>
		// 														<body>
		// 															<h1>Internal Server Error</h1>
		// 															<p>Okay, you know what? This one is on me.</p>
		// 														</body>
		// 													</html>`, statusCode)
	}
	w.writerState = stateHeader
	w.Writer.Write([]byte(statusLine))
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var str string
	for key := range headers {
		str += key + ": " + headers[key] + "\r\n"
	}
	str += "\r\n"

	w.Writer.Write([]byte(str))

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriterChunkedBody(p []byte) (int, error) {

}

func (w *Writer) WriterChunkedBodyDone() (int, error) {

}


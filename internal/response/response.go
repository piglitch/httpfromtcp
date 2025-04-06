package response

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"

	"go.uber.org/zap/buffer"
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
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

type Writer struct{
	Write 				*io.Writer
	writerState writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerState != stateStatus {
		return errors.New("Have to write status first")	
	}

	reasonHtml := ""
	if statusCode == 200 {
		statusCode = StatusOk
		reasonHtml = fmt.Sprintf(`<html>
																<head>
																	<title>%d OK</title>
																</head>
																<body>
																	<h1>Success!</h1>
																	<p>Your request was an absolute banger.</p>
																</body>
															</html>`, statusCode)
	}

	if statusCode == 400 {
		statusCode = StatusBadRequest
		reasonHtml = fmt.Sprintf(`<html>
										<head>
											<title>%d Bad Request</title>
										</head>
										<body>
											<h1>Bad Request</h1>
											<p>Your request honestly kinda sucked.</p>
										</body>
									</html>`, statusCode)
	}

	if statusCode == 500 {
		statusCode = StatusInternalError
		reasonHtml = fmt.Sprintf(`<html>
																<head>
																	<title>%d Internal Server Error</title>
																</head>
																<body>
																	<h1>Internal Server Error</h1>
																	<p>Okay, you know what? This one is on me.</p>
																</body>
															</html>`, statusCode)
	}
	fmt.Fprint(*w.Write, []byte(reasonHtml))
	w.writerState = stateHeader
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var str string
	for key := range headers {
		str += key + ": " + headers[key] + "\r\n"
	}
	str += "\r\n"
	fmt.Fprint(*w.Write, []byte(str))

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	
}
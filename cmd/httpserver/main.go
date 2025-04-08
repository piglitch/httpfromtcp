package main

import (
	// "crypto/sha256"
	"crypto/sha256"
	"fmt"

	// "httpfromtcp/internal/headers"
	// "httpfromtcp/internal/headers"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func ProxyHandler(writer response.Writer, trimmedPath string, req request.Request) {
	buff := make([]byte, 1024)
		
	resp, err := http.Get("https://httpbin.org/" + trimmedPath)
	if err != nil {
		return 
	}
	writer.WriteStatusLine(response.StatusOk)
	req.Headers.RemoveHeaders("Content-Length")
	req.Headers.Set("Transfer-Encoding", "chunked")
	req.Headers.Set("Content-Type", "text/html")
	
	trailers := []string{"X-Content-SHA256", "X-Content-Length"}
	trailersStr := trailers[0] + ", " + trailers[1]
	req.Headers.Set("Trailer", trailersStr)

	writer.WriteHeaders(req.Headers)
	for {
		n, err := resp.Body.Read(buff)
		if n > 0 {
			writer.WriterChunkedBody(buff[:n])
		}
		if err != nil {
			break
		}
		println(n)
	}
	// writer.WriteTrailers(trailerHeaders)
	writer.WriterChunkedBodyDone()

	trailerHeaders := headers.Headers{}
	hash := sha256.Sum256(writer.FullBody)
	trailerHeaders.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
	trailerHeaders.Set("X-Content-Length", fmt.Sprintf("%d", len(writer.FullBody)))
	// fmt.Printf("line 56: %d", fmt.Sprintf(hash))
	writer.WriteTrailers(trailerHeaders)
}

func HandlerFunc(w response.Writer, req *request.Request) {
	body := fmt.Sprintf(`<html>
													<head>
														<title>%d OK</title>
													</head>
													<body>
														<h1>Success!</h1>
														<p>Your request was an absolute banger.</p>
													</body>
												</html>`, 200)

	if req.RequestLine.RequestTarget == "/yourproblem"{
		w.WriteStatusLine(response.StatusBadRequest)
		body = fmt.Sprintf(`<html>
										<head>
											<title>%d Bad Request</title>
										</head>
										<body>
											<h1>Bad Request</h1>
											<p>Your request honestly kinda sucked.</p>
										</body>
									</html>`, response.StatusBadRequest)
	}

	if req.RequestLine.RequestTarget == "/myproblem"{
		w.WriteStatusLine(response.StatusInternalError)
		body = fmt.Sprintf(`<html>
													<head>
														<title>%d Internal Server Error</title>
													</head>
													<body>
														<h1>Internal Server Error</h1>
														<p>Okay, you know what? This one is on me.</p>
													</body>
												</html>`, response.StatusInternalError)
	}
	
	isHttpBinRoute := strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	
	if isHttpBinRoute {
		httpBinRoute := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
		ProxyHandler(w, httpBinRoute, *req)
		return
	}

	w.WriteStatusLine(response.StatusOk)

	// req.Headers.RemoveHeaders("Content-Length")
	w.WriteHeaders(req.Headers)
	
	w.WriteBody([]byte(body))
}

func main() {
	server, err := server.Serve(port, HandlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
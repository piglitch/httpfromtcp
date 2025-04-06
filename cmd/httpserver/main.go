package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

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

	w.WriteStatusLine(response.StatusOk)

	req.Headers.Set("Content-Type", "text/html")
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
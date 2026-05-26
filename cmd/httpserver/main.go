package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/calveshelder/httpfromlayer4/internal/request"
	"github.com/calveshelder/httpfromlayer4/internal/response"
	"github.com/calveshelder/httpfromlayer4/internal/server"
)

const port = 42069

func myHandler(w *response.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{Status: 400, Message: `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`}
	case "/myproblem":
		return &server.HandlerError{Status: 500, Message: `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`}
	default:
		body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
		w.WriteStatusLine(response.StatusOk)
		headers := response.GetDefaultHeaders(len(body))
		w.WriteHeaders(headers)
		w.Write([]byte(body))
	}
	return nil
}

func main() {
	server, err := server.Serve(port, myHandler)
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

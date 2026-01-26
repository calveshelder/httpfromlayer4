package main

import (
	"fmt"
	"github.com/calveshelder/httpfromlayer4/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		// Wait a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		rl, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal("Invalid request")
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", rl.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rl.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", rl.RequestLine.HttpVersion)
	}

}

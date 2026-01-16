package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	lines := getLinesChannel(f)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

// I don't fully understand this entire signature, what is <-chan doing here, is the function a receiver to chan or is it returning chan string (unlikely).
// Update:
// Got it. Function takes a parameter f of io.ReadCloser and will return a receive-only channel for strings to a caller.
func getLinesChannel(f io.ReadCloser) <-chan string {
	// Channel of strings.
	chanStrings := make(chan string)
	currentLine := ""
	b := make([]byte, 8)

	go func() {
		defer f.Close()
		defer close(chanStrings)

		for {
			n, err := f.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
				return

			}
			// After reading 8 bytes, split the data on newlines (\n) to create a slice of strings - let's call these split sections "parts". There will typically only be one or two "parts" because we're only reading 8 bytes at a time.
			str := string(b[:n])
			parts := []string{}
			parts = strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLine = currentLine + parts[i]
				// Send line to channel.
				chanStrings <- currentLine
				currentLine = ""
			}
			currentLine = currentLine + parts[len(parts)-1]
		}

		if currentLine != "" {
			chanStrings <- currentLine
		}
	}()

	return chanStrings
}

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	getLinesChannel(f)

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	// Channel of strings.
	chan_strings := make(chan string)
	current_line := ""

	go func() {
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			// After reading 8 bytes, split the data on newlines (\n) to create a slice of strings - let's call these split sections "parts". There will typically only be one or two "parts" because we're only reading 8 bytes at a time.
			str := string(b[:n])
			parts := []string{}
			parts = strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				current_line = current_line + parts[i]
				// Send line to channel.
				chan_strings <- current_line
				current_line = ""
			}
			current_line = current_line + parts[len(parts)-1]
		}
	}()

	if current_line != "" {
		fmt.Printf("read: %s\n", current_line)
	}

	return nil
}

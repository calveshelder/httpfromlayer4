package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		str := string(b[:n])
		fmt.Printf("read: %s\n", str)
	}
}

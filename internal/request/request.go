package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	requestStateInitialised = iota
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 8, 8)
	readToIndex := 0

	r := &Request{
		state: requestStateInitialised,
	}

	for r.state != requestStateDone {
		// If buffer is full, grow it
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read into buffer starting at ____
		n, err := reader.Read(buf[readToIndex:])

		if err == io.EOF {
			// What do you do here?
			r.state = requestStateDone
			break
		}

		readToIndex += n // Update our tracking variable

		// Parse what we have so far
		parsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Remove parsed data from buffer
		// Hint: copy(destination, source)
		copy(buf, buf[parsed:readToIndex])

		readToIndex -= parsed // Adjust tracking variable
	}

	return r, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	lineEnd := strings.Index(string(data), "\r\n")

	if lineEnd == -1 {
		return RequestLine{}, 0, nil
	}

	lineBytes := data[:lineEnd]
	line := string(lineBytes)
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid string")
	}

	bytesConsumed := lineEnd + 2

	method := parts[0]
	target := parts[1]
	versionNumber := strings.Split(parts[2], "/")
	number := versionNumber[1]

	for _, v := range parts[0] {
		isLetter := unicode.IsLetter(v)
		isUpper := unicode.IsUpper(v)
		isValid := isLetter && isUpper

		if !isValid {
			return RequestLine{}, 0, errors.New("Method is not valid")
		}
	}

	rl := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   number,
	}

	return rl, bytesConsumed, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == requestStateInitialised {
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return n, nil
		}
		r.RequestLine = rl
		r.state = requestStateDone
		return n, nil
	}
	if r.state == requestStateDone {
		return 0, errors.New("Trying to read data in a done state")
	}

	return 0, errors.New("unknown state")

}

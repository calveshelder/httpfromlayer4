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
	Headers     map[string]string
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	requestStateInitialised = iota
	requestStateParsingHeaders
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	r := &Request{
		state:   requestStateInitialised,
		Headers: make(map[string]string),
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
	totalBytesParsed := 0

	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}

		if n == 0 {
			return n, nil
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {
	// Switch/case logic
	switch r.state {
	case requestStateInitialised:
		//parseRequestLine
		rl, n, err := parseRequestLine(data)

		if err != nil {
			return n, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = rl
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		k, v, n, err := parseHeader(data)
		if err != nil {
			return n, err
		}

		if k == "" && v == "" {
			r.state = requestStateDone
			return n, nil
		}

		r.Headers[k] = v
		return n, err

	case requestStateDone:
		return 0, errors.New("trying to parse in done state")
	default:
		return 0, errors.New("unknown state")
	}

	// What do we return outside?
	return 0, errors.New("unreachable")
}

func parseHeader(data []byte) (key string, value string, bytesConsumed int, err error) {
	// 1. Find \r\n
	lineEnd := strings.Index(string(data), "\r\n")

	// 2. If not found, need more data.
	if lineEnd == -1 {
		return "", "", 0, nil
	}

	if lineEnd == 0 {
		return "", "", 2, nil
	}

	// 3. Split by the colon to get the key and the value.
	line := string(data[:lineEnd])
	split := strings.SplitN(line, ":", 2)

	// 4. Check if we got exactly 2 parts (key and value).
	if len(split) != 2 {
		return "", "", 0, errors.New("Malformed header: missing colon")
	}

	key = split[0]
	value = split[1]

	// 5. Lowecase key, trim whitespace from value.
	key = strings.ToLower(key)
	value = strings.TrimSpace(value)

	// 6. Calculate bytes consumed.
	bytesConsumed = lineEnd + 2 // + 2 is for \r\n.

	return key, value, 0, nil
}

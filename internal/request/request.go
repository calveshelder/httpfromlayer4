package request

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	requestString := string(request)
	requestStringParts := strings.Split(requestString, "\r\n")
	firstLine := requestStringParts[0]
	parsedRequestLine, err := parseRequestLine(firstLine)
	if err != nil {
		return nil, err
	}

	req := &Request{
		RequestLine: parsedRequestLine,
	}

	return req, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	parts := strings.Split(firstLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid string")
	}

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

	return rl, 0, nil
}

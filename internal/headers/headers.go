package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	lineEnd := strings.Index(string(data), "\r\n")
	if lineEnd == -1 {
		return 0, false, nil
	}
	if lineEnd == 0 {
		return 2, true, nil
	}
	rawLine := string(data[:lineEnd])
	trimmed := strings.TrimSpace(rawLine)
	colonIndex := strings.Index(trimmed, ":")

	if colonIndex == -1 {
		return 0, false, errors.New("invalid header")
	}

	if colonIndex > 0 && trimmed[colonIndex-1] == ' ' {
		return 0, false, errors.New("invalid header")
	}

	key := strings.TrimSpace(trimmed[:colonIndex])
	value := strings.TrimSpace(trimmed[colonIndex+1:])

	h[key] = value

	n = lineEnd + len("\r\n")

	return n, false, nil
}

package response

import (
	"fmt"
	"io"

	"github.com/calveshelder/httpfromlayer4/internal/headers"
)

type StatusCode int

const (
	StatusOk StatusCode = iota
	BadRequest
	InternalServerError
)

var statusName = map[StatusCode]string{
	StatusOk:            "200 OK",
	BadRequest:          "400 Bad Request",
	InternalServerError: "500 Internal Server Error",
}

func (ss StatusCode) String() string {
	return statusName[ss]
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := fmt.Fprintf(w, "HTTP/1.1 %s\r\n", statusCode.String())
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\r\n")

	return err

}

package response

import (
	"fmt"
	"io"

	"github.com/calveshelder/httpfromlayer4/internal/headers"
)

type StatusCode int

type Writer struct {
	writer io.Writer
}

const (
	StatusOk            StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
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
	h["Content-Type"] = "text/html"
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

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return WriteStatusLine(w.writer, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w.writer, headers)
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

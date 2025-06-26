package server

import (
	"io"

	"github.com/mharkness1/httpfromtcp/internal/request"
	"github.com/mharkness1/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type HandlerFunc func(w io.Writer, req *request.Request) *HandlerError

func (err HandlerError) HandlerErrorWriter(w io.Writer) {
	response.WriteStatusLine(w, err.StatusCode)
	message := []byte(err.Message)
	headers := response.GetDefaultHeaders(len(message))
	response.WriteHeaders(w, headers)
	w.Write(message)
}

package server

import (
	"github.com/mharkness1/httpfromtcp/internal/request"
	"github.com/mharkness1/httpfromtcp/internal/response"
)

/*
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
*/

type HandlerFunc func(w *response.Writer, req *request.Request)

/*
func (err HandlerError) HandlerErrorWriter(w io.Writer) {
	response.WriteStatusLine(w, err.StatusCode)
	message := []byte(err.Message)
	headers := response.GetDefaultHeaders(len(message))
	response.WriteHeaders(w, headers)
	w.Write(message)
}
*/

package response

import (
	"fmt"

	"github.com/mharkness1/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writerStateTrailers {
		return fmt.Errorf("writer in wrong state: %d", w.writerState)
	}
	for k, v := range h {
		line := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error writing trailer line: %v\n", err)
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		fmt.Printf("Error writing final CRLF: %v\n", err)
		return err
	}
	w.writerState = writerStateBody
	fmt.Printf("WriteTrailers completed successfully\n")
	return nil
}

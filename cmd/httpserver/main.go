package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mharkness1/httpfromtcp/internal/headers"
	"github.com/mharkness1/httpfromtcp/internal/request"
	"github.com/mharkness1/httpfromtcp/internal/response"
	"github.com/mharkness1/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(http400)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(http500)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(http200)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeSuccess)
	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-Sha256, X-Content-Length")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	fullBody := make([]byte, 0)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("error writing chunked body:", err)
				break
			}
			fullBody = append(fullBody, buffer[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading resp body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body end:", err)
	}
	trailers := headers.NewHeaders()

	sha256Hash := fmt.Sprintf("%x", sha256.Sum256(fullBody))

	trailers["X-Content-Sha256"] = sha256Hash
	trailers["X-Content-Length"] = fmt.Sprintf("%d", len(fullBody))

	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Printf("Error writing trailers: %v\n", err)
	} else {
		fmt.Printf("WriteTrailers returned successfully\n")
	}
}

const http400 string = `<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`

const http500 string = `<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`

const http200 string = `<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`

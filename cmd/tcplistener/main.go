package main

import (
	"fmt"
	"net"

	"github.com/mharkness1/httpfromtcp/internal/request"
)

const tcpPort = ":42069"

func main() {
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Printf("error opening connection: %v", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: %v", err)
			break
		}

		fmt.Print("Network connection accepted.\n")
		req, err := request.RequestFromReader(conn)

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Printf("Body:")
		fmt.Printf("%s\n", req.Body)

		fmt.Print("Connection closed.")
	}
}

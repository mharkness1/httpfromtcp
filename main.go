package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
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
		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Printf("%v\n", line)
		}
		fmt.Print("Connection closed.")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error %v", err.Error())
				return
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}

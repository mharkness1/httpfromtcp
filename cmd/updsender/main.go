package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const udpPort = ":42069"

func main() {
	updAddr, err := net.ResolveUDPAddr("udp", "localhost"+udpPort)
	if err != nil {
		fmt.Printf("error resolving address: %v", err)
		return
	}

	updConn, err := net.DialUDP("udp", nil, updAddr)
	if err != nil {
		fmt.Printf("error creating connection: %v", err)
		return
	}
	fmt.Println("Connection established UDP")
	defer updConn.Close()

	readerPtr := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n>")
		line, err := readerPtr.ReadString('\n')
		if err != nil {
			fmt.Printf("error getting line: %v", err)
		}
		_, err = updConn.Write([]byte(line))
		if err != nil {
			fmt.Printf("error writing to connection: %v", err)
		}
	}
}

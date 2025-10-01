package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// const address = "127.0.0.1:42069"
const address = "localhost:42069"
const protocol = "udp"

func main() {
	udpAddress, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		log.Fatalf("Error resolving UDP address %s: %s", address, err)
	}

	udpConn, err := net.DialUDP(protocol, nil, udpAddress)
	if err != nil {
		log.Fatalf("Error creating UDP conn %s: %s", udpAddress, err)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString(byte('\n'))
		if err != nil {
			log.Fatalf("Error reading from buffer: %s", err)
		}

		_, err = udpConn.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error writing to UDP conn: %s", err)
		}
	}
}

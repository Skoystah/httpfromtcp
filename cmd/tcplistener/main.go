package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const address = "127.0.0.1:42069"
const protocol = "tcp"

func main() {

	listener, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("Error creating listener via %s on %s", protocol, address)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection via %s on %s", protocol, address)
		}

		fmt.Printf("===Connection has been accepted via %s on %s===\n", protocol, address)
		fmt.Printf("=======================================\n")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error reading request: %s", err)
		}

		fmt.Print("Request line:\n")
		fmt.Printf("- Method: %s", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Print("Headers:\n")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Print("Body:\n")
		fmt.Printf("%s\n", string(request.Body))
		fmt.Printf("===Connection has been closed===\n")
	}
}

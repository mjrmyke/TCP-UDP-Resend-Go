package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server starting up")
	go UDPServer("8081")
	go TCPServer("8082")

	for {

	}
}

func TCPServer(port string) {

	fmt.Println("in TCP Server")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up TCPServer on Port: ", port)
	}

	connection, err := ln.Accept()
	if err != nil {
		fmt.Println("error Accepting in TCPServer on Port: ", port)
	}

	for {
		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println("error reading in a string ")
			fmt.Print("Message Received:", string(message))

		}

	}

}

func UDPServer(port string) {

	fmt.Println("in UDP Server")
	ln, err := net.Listen("UDP", port)
	if err != nil {
		fmt.Println("error starting up UDPServer on Port: ", port)
	}

	connection, err := ln.Accept()
	if err != nil {
		fmt.Println("error Accepting in UDPServer on Port: ", port)
	}

	for {
		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println("error reading in a string ")
			fmt.Print("Message Received:", string(message))

		}

	}

}

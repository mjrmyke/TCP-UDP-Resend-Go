package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server starting up")
	go UDPServer(":8085")
	go TCPServer(":8082")

	for {

	}
}

//TCPServer is a function to handle incoming TCP connections
func TCPServer(port string) {

	fmt.Println("in TCP Server")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up TCPServer on Port: ", port, err)
	}
	fmt.Println("waiting for TCP INFO")

	connection, err := ln.Accept()
	if err != nil {
		fmt.Println("error Accepting in TCPServer on Port: ", port)
	}

	for {
		fmt.Println("receiving tcp input")

		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println("error reading in a string ")
			fmt.Print("Message Received:", string(message))

		}

	}

}

//UDPServer is a function to handle incoming TCP connections
func UDPServer(port string) {

	fmt.Println("in UDP Server")
	serveraddr, err := net.ResolveUDPAddr("udp", port)

	ln, err := net.ListenUDP("udp", serveraddr)
	if err != nil {
		fmt.Println("error starting up UDPServer on Port: ", port, err)
	}

	buf := make([]byte, 1024)

	for {
		fmt.Println("waiting for UDP INFO")

		n, addr, err := ln.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)

		}
	}

}

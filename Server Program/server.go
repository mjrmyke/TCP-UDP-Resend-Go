package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

//Packet struct is a datatype to hold a header number to keep track of operations
//and Data is the string sent
type Packet struct {
	HeaderNum int    `json:"headernum"`
	Data      string `json:"data"`
}

var mainwg sync.WaitGroup
var tcpwg sync.WaitGroup
var udpwg sync.WaitGroup

func main() {
	fmt.Println("Server starting up")
	go UDPServer(":8085")
	go TCPServer(":8082")
	mainwg.Add(2)
	mainwg.Wait()

}

//TCPServer is a function to handle incoming TCP connections
func TCPServer(port string) {

	fmt.Println("in TCP Server")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up TCPServer on Port: ", port, err)
	}
	fmt.Println("waiting for TCP INFO")

	//sets the last the program will do before exiting, close the connections
	defer ln.Close()

	connection, err := ln.Accept()
	if err != nil {
		fmt.Println("error Accepting in TCPServer on Port: ", port)
	}
	go HandleIncomingTCPData(connection)
	tcpwg.Add(1)
	tcpwg.Wait()

	defer mainwg.Done()

}

//HandleIncomingTCPData is a function to deal with the incoming TCP connection
func HandleIncomingTCPData(connection net.Conn) {
	// message := ""

	for {

		message, _ := bufio.NewReader(connection).ReadString('\n')
		// if err != nil {
		// 	break
		// }
		if message != "" {
			fmt.Print("Message Received:", string(message))
			connection.Write([]byte("Received, and redirecting:" + message))
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

	//sets the last the program will do before exiting, close the connections
	defer ln.Close()

	buf := make([]byte, 1024)

	for {
		fmt.Println("waiting for UDP INFO")

		n, addr, err := ln.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)

		}
	}
	defer mainwg.Done()

}

package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	fmt.Println("Entering Client Program for CSCI 156 Project")
	var TCPIP, UDPIP string
	if len(os.Args) == 1 {
		TCPIP = "127.0.0.1:8082"
		UDPIP = "127.0.0.1:8085"
	} else {
		TCPIP = os.Args[1]
		UDPIP = os.Args[2]
	}

	fmt.Println("TCPIP IS: ", TCPIP)
	fmt.Println("UDPIP IS: ", UDPIP)
	var wg sync.WaitGroup

	go TCPConnect(TCPIP)
	go UDPConnect(UDPIP)
	wg.Add(2)
	wg.Wait()
	// tcpconn, err := net.Dial("tcp", TCPIP)
	// if err != nil {
	// 	log.Println("error while dialing TCP", err, TCPIP)
	// }
	// tcpconn.Write(makedataforpacket("HELLO\n"))

	// err = tcpconn.Close()
	// if err != nil {
	// 	log.Println("Closing Connection", err, TCPIP)
	// }
}

//TCPConnect is a function that creates a TCP connection, and sends a small message
func TCPConnect(TCPIP string) {
	tcpconn, err := net.Dial("tcp", TCPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, TCPIP)
	}
	tcpconn.Write(makedataforpacket("HELLO TCP!\n"))
	fmt.Println("sent data in TCP")
	err = tcpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, TCPIP)
	}
}

//UDPConnect is a function that creates a TCP connection, and sends a small message
func UDPConnect(UDPIP string) {
	udpconn, err := net.Dial("udp", UDPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, UDPIP)
	}
	udpconn.Write(makedataforpacket("HELLO UDP!\n"))
	fmt.Println("sent data in UDP")

	err = udpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, UDPIP)
	}
}

//function makedataforpacket receives an integer index number,
//converts that number to bytes, then places it in a
//1000 byte buffer and returns it
func makedataforpacket(index string) []byte {
	//convert int to slice of bytes
	indexdata := []byte(index)
	//place that slice of bytes into a buffer
	bytebuffer := bytes.NewBuffer(indexdata)
	return bytebuffer.Bytes()

}

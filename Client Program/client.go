package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

//Packet struct is a datatype to hold a header number to keep track of operations
//and Data is the string sent
type Packet struct {
	HeaderNum int    `json:"headernum"`
	Data      string `json:"data"`
}

var wg sync.WaitGroup
var MAXPACKETS = 10000

func main() {
	fmt.Println("Entering Client Program for CSCI 156 Project")
	var TCPIP, UDPIP, IP string
	if len(os.Args) == 1 {
		IP = "127.0.0.1"
	} else {
		IP = os.Args[1]
	}
	TCPIP = IP + ":8082"
	UDPIP = IP + ":8085"

	fmt.Println("TCPIP IS: ", TCPIP)
	fmt.Println("UDPIP IS: ", UDPIP)

	go TCPConnect(TCPIP)
	go UDPConnect(UDPIP)
	wg.Add(2)
	wg.Wait()

}

//TCPConnect is a function that creates a TCP connection, and sends a small message
func TCPConnect(TCPIP string) {
	tcpconn, err := net.Dial("tcp", TCPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, TCPIP)
	}

	for k := 0; k < MAXPACKETS; k++ {
		// for k := 0; k < 100; k++ {

		tcpconn.Write(makedataforpacket("obj number: " + strconv.FormatInt(int64(k), 10) + " HELLO TCP!\n"))
		fmt.Println("sent data in TCP")
		message, _ := bufio.NewReader(tcpconn).ReadString('\n')
		fmt.Println("Received: ", message)

		// buf := make([]byte, 1024)
		// n, addr, err := udpconn.ReadFromUDP(buf)
		// if err != nil {
		// 	log.Println("error reading UDP", err, UDPIP)
		// }
		// fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		time.Sleep(5 * time.Millisecond)

	}

	err = tcpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, TCPIP)
	}
	defer wg.Done()
}

//UDPConnect is a function that creates a TCP connection, and sends a small message
func UDPConnect(UDPIP string) {
	udpconn, err := net.Dial("udp", UDPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, UDPIP)
	}
	udpconn.Write(makedataforpacket("HELLO UDP!\n"))
	fmt.Println("sent data in UDP")

	for k := 0; k < MAXPACKETS; k++ {

		udpconn.Write(makedataforpacket("obj number: " + strconv.FormatInt(int64(k), 10) + "HELLO UDP!\n"))
		time.Sleep(5 * time.Millisecond)

	}

	err = udpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, UDPIP)
	}
	defer wg.Done()
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

// buf := make([]byte, 1024)
// n, addr, err := udpconn.ReadFromUDP(buf)
// if err != nil {
// 	log.Println("error reading UDP", err, UDPIP)
// }
// fmt.Println("Received ", string(buf[0:n]), " from ", addr)

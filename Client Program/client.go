package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

//Packet struct is a datatype to hold a header number to keep track of operations
//and Data is the string sent
type Packet struct {
	conn      net.Conn
	ConnType  string
	HeaderNum int    `json:"headernum"`
	Data      string `json:"data"`
}

type LostPackets struct {
	Packets []int `json:"packets"`
}

var resendlist = make(chan Packet)
var wg sync.WaitGroup
var MAXPACKETS = 7000
var numretry, packetsreceived []int32

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
	// wg.Add(1)
	wg.Wait()

}

//TCPConnect is a function that creates a TCP connection, and sends a small message
func TCPConnect(TCPIP string) {
	tcpconn, err := net.Dial("tcp", TCPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, TCPIP)
	}

	fmt.Println("Initiating TCP Connection to Server")
	for {

		// message, err := bufio.NewReader(tcpconn).ReadString('\n')
		// if err != nil {
		// 	log.Println("error message, err := bufio.NewReader(tcpconn).ReadString('\n')", err)
		// }

		// fmt.Println("Received: ", message)

		// // buf := make([]byte, 1024)
		// // n, addr, err := udpconn.ReadFromUDP(buf)
		// // if err != nil {
		// // 	log.Println("error reading UDP", err, UDPIP)
		// // }
		// // fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		decoder := json.NewDecoder(tcpconn)
		var tmplist LostPackets

		err := decoder.Decode(&tmplist)
		if err != nil {
			fmt.Printf("Client disconnected.\n")
			break
		}
		for _, k := range tmplist.Packets {
			tmppacket := Packet{
				HeaderNum: k,
				Data:      "this is misc packetdata",
				ConnType:  "UDP",
			}

			resendlist <- tmppacket

		}

	}

	err = tcpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, TCPIP)
	}
	defer wg.Done()
}

func resendpackets(conn net.Conn) {
	for {
		//wait for data to resend
		tmppacket := <-resendlist

		tmpdata, err := json.Marshal(&tmppacket)
		if err != nil {
			log.Println("error Json.Marshal(tmppacket)", err)
		}

		_, err = conn.Write(tmpdata)
		fmt.Println("Resending Packet num:", tmppacket)
		if err != nil {
			log.Println("error conn.Write(tmpdata)", err)
		}
		time.Sleep(3 * time.Millisecond)
	}

}

//UDPConnect is a function that creates a TCP connection, and sends a small message
func UDPConnect(UDPIP string) {
	udpconn, err := net.Dial("udp", UDPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, UDPIP)
	}
	fmt.Println("sent data in UDP")
	go resendpackets(udpconn)

	for k := 0; k < MAXPACKETS; k++ {
		tmppacket := &Packet{
			HeaderNum: k,
			Data:      "this is misc packetdata",
			ConnType:  "UDP",
		}

		tmpdata, err := json.Marshal(tmppacket)
		if err != nil {
			log.Println("error Json.Marshal(tmppacket)", err)
		}

		if rand.Intn(100) > 90 {
			fmt.Println("simulated packet loss")
			//do nothing
		} else {
			udpconn.Write(makedataforpacket(tmpdata))
			fmt.Println("sent data in UDP")
		}

		time.Sleep(1 * time.Millisecond)

	}

	// err = udpconn.Close()
	// if err != nil {
	// 	log.Println("Closing Connection", err, UDPIP)
	// }
	defer wg.Done()
}

//function makedataforpacket receives an integer index number,
//converts that number to bytes, then places it in a
//1000 byte buffer and returns it
func makedataforpacket(data []byte) []byte {

	//place that slice of bytes into a buffer
	bytebuffer := bytes.NewBuffer(data)
	return bytebuffer.Bytes()

}

// buf := make([]byte, 1024)
// n, addr, err := udpconn.ReadFromUDP(buf)
// if err != nil {
// 	log.Println("error reading UDP", err, UDPIP)
// }
// fmt.Println("Received ", string(buf[0:n]), " from ", addr)

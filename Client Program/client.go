////////////////////////////////////////////////////////////////////////////////////////////////////
// Myke Walker CSCI 156 Project
// Sending and receiving Data through UDP & TCP
// This is the Client program, which sends up to MAXPACKETS to a server,
// simulates packet loss, and receives information on which packets that are needed to be resent
///////////////////////////////////////////////////////////////////////////////////////////////////
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

//LostPackets is a struct that is used to contain a slice of ints,
//which indicate what packets will need to be resent
type LostPackets struct {
	Packets []int `json:"packets"`
}

//variables for statistics
var lostpackets int
var sentpackets int

//pass percent is the percent that the packet will go through successfully
var passpercent = 90

//channel to pass the list of packets to resend
var resendlist = make(chan Packet)

//waitgroup to sync threads
var wg sync.WaitGroup

//MAXPACKETS of packets to be sent
var MAXPACKETS = 10000

//slice of ints for amount to retry, and packetsreceived
var numretry, packetsreceived []int32

func main() {
	fmt.Println("Entering Client Program for CSCI 156 Project")
	//declar elocal vars
	var TCPIP, UDPIP, IP string
	//get the ip address to send the information to if passed, or use default
	if len(os.Args) == 1 {
		IP = "127.0.0.1"
	} else {
		IP = os.Args[1]
	}
	//add ports to the ips
	TCPIP = IP + ":8082"
	UDPIP = IP + ":8085"

	fmt.Println("TCPIP IS: ", TCPIP)
	fmt.Println("UDPIP IS: ", UDPIP)

	//start up servers on their own threads
	go TCPConnect(TCPIP)
	go UDPConnect(UDPIP)
	go printstats()
	//have main program wait for the jobs to finish
	wg.Add(2)
	wg.Wait()

}

func printstats() {
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Number of Failed Sends:", lostpackets)
		fmt.Println("Number of Sends:", sentpackets)
	}
}

//TCPConnect is a function that creates a TCP connection
func TCPConnect(TCPIP string) {
	//dial server at specified ip
	tcpconn, err := net.Dial("tcp", TCPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, TCPIP)
	}

	fmt.Println("Initiated TCP Connection to Server")
	//wait for server to send the status of the data transfer
	for {

		//add decoder to the connection information
		decoder := json.NewDecoder(tcpconn)
		var tmplist LostPackets

		//decode the information sent
		err := decoder.Decode(&tmplist)
		if err != nil {
			fmt.Printf("Client disconnected.\n", err)
			break
		}
		//create each individual packet to be resent, then pass it into the channel for the other threads
		for _, k := range tmplist.Packets {
			tmppacket := Packet{
				HeaderNum: k,
				Data:      "this is misc packetdata",
				ConnType:  "UDP",
			}

			//pass the packet to another thread
			resendlist <- tmppacket

		}

		time.Sleep(5 * time.Millisecond)

	}

	err = tcpconn.Close()
	if err != nil {
		log.Println("Closing Connection", err, TCPIP)
	}
	defer wg.Done()
}

//resendpackets is a function that runs on its own thread,
//and waits for packets to be resent from the channel.
//when it gets one, its ends it and waits 3 milliseconds
func resendpackets(conn net.Conn) {

	for {
		//wait for data from channel (other threads)
		tmppacket := <-resendlist

		//marshal packet into json
		tmpdata, err := json.Marshal(&tmppacket)
		if err != nil {
			log.Println("error Json.Marshal(tmppacket)", err)
		}

		if rand.Intn(100) > passpercent {
			//do nothing
			fmt.Println("Simulated packet loss on a packet to be resent")
			lostpackets++
		} else {
			sentpackets++
			//write the marshalled json to the connection
			_, err = conn.Write(tmpdata)
			fmt.Println("Resending Packet num:", tmppacket)
			if err != nil {
				log.Println("error conn.Write(tmpdata)", err)
			}
		}
		//wait for 3 ms
		time.Sleep(3 * time.Millisecond)
	}

}

//UDPConnect is a function that creates a udp connection, and an amount of data
//specified from MAXPACKETS
func UDPConnect(UDPIP string) {
	//dial a udp server
	udpconn, err := net.Dial("udp", UDPIP)
	if err != nil {
		log.Println("error while dialing TCP", err, UDPIP)
	}
	//start up the goroutine to monitor which packets need to be resent
	go resendpackets(udpconn)

	//loop from 0 to MAXPACKETS
	for k := 0; k < MAXPACKETS; k++ {
		//creates each packet through the loop
		tmppacket := &Packet{
			HeaderNum: k,
			Data:      "this is misc packetdata",
			ConnType:  "UDP",
		}

		//marshal up the datatype
		tmpdata, err := json.Marshal(tmppacket)
		if err != nil {
			log.Println("error Json.Marshal(tmppacket)", err)
		}

		//simulate packet loss
		if rand.Intn(100) > passpercent {
			fmt.Println("simulated packet loss on packet num: ", k)
			lostpackets++

			//do nothing
		} else {
			sentpackets++
			//write the data
			udpconn.Write(makedataforpacket(tmpdata))
			fmt.Println("sent data in UDP, packet num: ", k)
		}
		//sleep 3 ms
		time.Sleep(3 * time.Millisecond)

	}

	//when program finishes, let the waitgroup know
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type LostPackets struct {
	Packets []int `json:"packets"`
}

//Packet struct is a datatype to hold a header number to keep track of operations
//and Data is the string sent
type UDPPacket struct {
	Addr      *net.UDPAddr
	Conn      net.Conn
	HeaderNum int    `json:"headernum"`
	Data      string `json:"data"`
}

var numpacketsrange = 20
var mainwg sync.WaitGroup
var tcpwg sync.WaitGroup
var udpwg sync.WaitGroup
var ackchan = make(chan int)
var resendchan = make(chan int)
var packets [10000]bool
var resendmap = make(map[int]bool, 11000)

func main() {
	go HandleIncomingACKs()
	fmt.Println("Server starting up")
	go UDPServer(":8085")
	go TCPServer(":8082")
	mainwg.Add(2)
	mainwg.Wait()

}

func StatusOfTransmission(conn net.Conn) {
	fmt.Println("Entered Status of Transmission Func, entering while loop")

	for {
		var resendlist []int
		fmt.Println("Number of missing packets", len(resendmap))
		fmt.Println("Missing Packets: ", resendmap)

		if len(resendmap) > 0 {

			for k := range resendmap {
				resendlist = append(resendlist, k)
			}

			tmplist := LostPackets{
				Packets: resendlist,
			}

			tmpdata, err := json.Marshal(tmplist)
			if err != nil {
				log.Println("error Json.Marshal(tmppacket)", err)
			}

			_, err = conn.Write(tmpdata)
			if err != nil {
				log.Println("error conn.Write(tmpdata)", err)
			}
		}

		time.Sleep(15 * time.Second)

	}

}

//TCPServer is a function to handle incoming TCP connections
func TCPServer(port string) {
	defer mainwg.Done()

	fmt.Println("in TCP Server")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up TCPServer on Port: ", port, err)
	}
	fmt.Println("waiting for TCP INFO")

	//sets the last the func will do before exiting, close the connections
	defer ln.Close()

	for {
		connection, err := ln.Accept()
		if err != nil {
			fmt.Println("error Accepting in TCPServer on Port: ", port)
			break
		}
		fmt.Println("Accepted connection", connection)
		go StatusOfTransmission(connection)

		go func() {

			for {
				//maintain open connection

			}
		}()

	}

}

//HandleIncomingTCPData is a function to deal with the incoming TCP connection
func HandleIncomingData(Packets chan UDPPacket) {
	fmt.Println("preparing to handle incoming tcp data, entering loop")
	for {
		// Wait for the next job to come off the queue.
		tmppacket := <-Packets
		// fmt.Println("got next packet off of the packet channel")

		// fmt.Println("Received:", tmppacket)

		// Send back the response.
		//prepare data for TCP ACK
		ackchan <- tmppacket.HeaderNum

	}

}

//HandleIncomingACKs is a function to deal with the incoming TCP connection
func HandleIncomingACKs() {
	fmt.Println("preparing to handle ACKS, entering loop")
	for {
		// Wait for the next job to come off the queue.
		tmpacknum := <-ackchan
		// fmt.Println("got ack from UDP Data")
		packets[tmpacknum] = true

		if resendmap[tmpacknum] == true {
			delete(resendmap, tmpacknum)
		}

		//every numpacketsrange packets, run the
		if tmpacknum%numpacketsrange == 0 {

			for j := 0; j <= 19; j++ {
				if (tmpacknum-j > 0) && (packets[tmpacknum-j] == false) {
					fmt.Println("Need Packet Num", tmpacknum-j)

					resendmap[tmpacknum-j] = true

				}

			}

		}
	}

}

//UDPServer is a function to handle incoming TCP connections
func UDPServer(port string) {
	defer mainwg.Done()
	Packets := make(chan UDPPacket)

	fmt.Println("in UDP Server")

	ln, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		fmt.Println("error starting up UDP on Port: ", port, err)
	}

	go HandleIncomingData(Packets)
	conn, err := net.ListenUDP("udp", ln)
	if err != nil {
		fmt.Println("error starting up UDP on Port: ", port, err)
	}
	udpwg.Add(1)
	err = conn.SetReadBuffer(10000)
	if err != nil {
		log.Println("	err = conn.SetReadBuffer(10000)", err)
	}
	err = conn.SetWriteBuffer(10000)
	if err != nil {
		log.Println("err = conn.SetReadBuffer(10000)", err)
	}

	go func() {
		buf := make([]byte, 1024)

		for {

			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				fmt.Printf("Client disconnected.\n")
				fmt.Println(err)
				break
			}
			var tmppacket UDPPacket

			err = json.Unmarshal(buf[0:n], &tmppacket)
			if err != nil {
				log.Print(err)
				return
			}
			tmppacket.Conn = conn
			tmppacket.Addr = ln

			Packets <- tmppacket

		}
	}()
	udpwg.Wait()

}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Myke Walker CSCI 156 Project
// Sending and receiving Data through UDP & TCP
// This is the Server program, which receives up to MAXPACKETS from a client,
// receives the packets, determines which packets werent sent, and requests for the packets to be resent
////////////////////////////////////////////////////////////////////////////////////////////////////////
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

//LostPackets is a struct that is used to contain a slice of ints,
//which indicate what packets will need to be resent
type LostPackets struct {
	//slice of packets
	Packets []int `json:"packets"`
}

//UDPPacket struct is a datatype to hold a header number to keep track of operations
//and Data is the string sent
type UDPPacket struct {
	//UDP address information
	Addr *net.UDPAddr
	//Actual UDP Connection
	Conn net.Conn
	//Number indicating which packet was sent
	HeaderNum int `json:"headernum"`
	//Misc data that will be take up at least 1000 bytes
	Data string `json:"data"`
}

//number of packets to look back for a "gobackn" type protocol
var numback = 50

//amount of times a packet will be retried
var retrylimit = 3

//var to track number of deleted packets from resendlist
var numdel int

//how many packets to look back (ie: GO back N Protocol)
var numpacketsrange = 20

//Waitgroups for various threads
var mainwg sync.WaitGroup
var tcpwg sync.WaitGroup
var udpwg sync.WaitGroup

//channels to pass information between threads
var ackchan = make(chan int)
var resendchan = make(chan int)

//Array of bools to indicate what has been sent
var packets [10000]bool

//Map of int to int, if an item exists in the map, it needs to be resent
var resendmap = make(map[int]int, 11000)

//entry point of application
func main() {
	//start goroutine to handle incoming ACKS
	go HandleIncomingACKs()
	fmt.Println("Server starting up")
	//start goroutines for each server
	go UDPServer(":8085")
	go TCPServer(":8082")
	// wait for the 2 servers to close
	mainwg.Add(2)
	mainwg.Wait()

}

//StatusOfTransmission is a func that determines monitors the resend map,
//if there are packets to be resent
func StatusOfTransmission(conn net.Conn) {
	fmt.Println("Entered Status of Transmission Func, entering while loop")
	counter := 0

	//while loop
	for {
		var resendlist []int
		fmt.Println("Number of missing packets", len(resendmap))
		fmt.Println("Missing Packets: ", resendmap)

		//if there exists packets to be sent
		if len(resendmap) > 0 {
			//reset wait counter
			counter = 0

			//make the list from the map
			for k := range resendmap {
				if resendmap[k] > retrylimit {
					fmt.Println("Gave up on retrieving packet: ", k)
					delete(resendmap, k)
					numdel++

				} else {
					resendlist = append(resendlist, k)
					resendmap[k]++
				}
			}

			//put the list in a datatype to be shipped
			tmplist := LostPackets{
				Packets: resendlist,
			}

			//struct to json
			tmpdata, err := json.Marshal(tmplist)
			if err != nil {
				log.Println("error Json.Marshal(tmppacket)", err)
			}

			//write json to connection
			_, err = conn.Write(tmpdata)
			if err != nil {
				log.Println("error conn.Write(tmpdata)", err)
			}
		}

		//if tried 5 times, and no more packets were needed
		//file is done being sent
		if counter > 5 {
			break
		}
		counter++
		//wait for 5 seconds, and then try again
		time.Sleep(5 * time.Second)

	}
	// fmt.Println("number of packets that were given up on: ", numdel)
	fmt.Println("File Transfer Completed")

}

//TCPServer is a function to handle incoming TCP connections
func TCPServer(port string) {
	defer mainwg.Done()

	fmt.Println("in TCP Server")
	//listen for TCP conns on the specified port
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up TCPServer on Port: ", port, err)
	}
	fmt.Println("waiting for TCP INFO")

	//sets the last the func will do before exiting, close the connections
	defer ln.Close()

	//for each connection, accept it, and then place it in its own thread
	for {
		connection, err := ln.Accept()
		if err != nil {
			fmt.Println("error Accepting in TCPServer on Port: ", port)
			break
		}
		fmt.Println("Accepted connection", connection)
		go StatusOfTransmission(connection)

	}

}

//HandleIncomingData is a function to deal with the incoming TCP connection
func HandleIncomingData(Packets chan UDPPacket) {
	fmt.Println("preparing to handle incoming tcp data, entering loop")
	for {
		// Wait for the next job to come off the queue.
		tmppacket := <-Packets

		//uncomment to debug received packets
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
		packets[tmpacknum] = true

		//if it is a packet that neeeded to be resent, delete it from the map
		if resendmap[tmpacknum] > 0 {
			delete(resendmap, tmpacknum)
		}

		//every numpacketsrange packets, run the
		if tmpacknum%numpacketsrange == 0 {

			//check the last numback packets before this one
			for j := 0; j <= numback; j++ {
				//if it is non negative, and i do not have the packet, mark it
				if (tmpacknum-j > 0) && (packets[tmpacknum-j] == false) {
					fmt.Println("Need Packet Num", tmpacknum-j)

					resendmap[tmpacknum-j] = 1

				}

			}

		}
	}

}

//UDPServer is a function to handle incoming TCP connections
func UDPServer(port string) {
	defer mainwg.Done()
	//make a channel to pass between goroutines, then pass it into one
	Packets := make(chan UDPPacket)
	go HandleIncomingData(Packets)

	fmt.Println("in UDP Server")

	//determine my information regarding a udp address
	ln, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		fmt.Println("error starting up UDP on Port: ", port, err)
	}

	//listen for udp datagrams on the udpaddress that was retreived
	conn, err := net.ListenUDP("udp", ln)
	if err != nil {
		fmt.Println("error starting up UDP on Port: ", port, err)
	}

	//make the buffer large enough for all packets (unnecessary)
	err = conn.SetReadBuffer(10000)
	if err != nil {
		log.Println("	err = conn.SetReadBuffer(10000)", err)
	}
	err = conn.SetWriteBuffer(10000)
	if err != nil {
		log.Println("err = conn.SetReadBuffer(10000)", err)
	}

	//anonymous function that runs on a gourtine to receive incoming packets
	go func() {
		buf := make([]byte, 1024)

		for {

			//read from the connection
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				fmt.Printf("Client disconnected.\n")
				fmt.Println(err)
				break
			}
			var tmppacket UDPPacket

			//unmarshal the json to a packet
			err = json.Unmarshal(buf[0:n], &tmppacket)
			if err != nil {
				log.Print(err)
				return
			}
			//pass along the connection information
			tmppacket.Conn = conn
			tmppacket.Addr = ln

			//pass the completed packet onto the channel
			Packets <- tmppacket

		}
	}()
	udpwg.Wait()

}

package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("Entering Client Program for CSCI 156 Project")

}

//function makedataforpacket receives an integer index number,
//converts that number to bytes, then places it in a
//1000 byte buffer and returns it
func makedataforpacket(index int) *bytes.Buffer {
	//convert int to slice of bytes
	indexdata := []byte(strconv.Itoa(index))
	//place that slice of bytes into a buffer
	bytebuffer := bytes.NewBuffer(indexdata)
	return bytebuffer

}

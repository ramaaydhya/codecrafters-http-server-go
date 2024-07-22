package main

import (
	"bytes"
	"fmt"
	"strconv"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	var req = make([]byte, 1024)
	conn.Read(req)

	if bytes.HasPrefix(req, []byte("GET / HTTP/1.1")) {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	var getEchoPrefix = []byte("GET /echo/")
	if !bytes.HasPrefix(req, getEchoPrefix) {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	var strBegin = len(getEchoPrefix) // GET /echo/
	var strEnd = strBegin
	for req[strEnd] != ' ' {
		strEnd++
	}

	var header = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: ")
	var conLen = strEnd - strBegin
	var bSlcConLen = []byte(strconv.Itoa(conLen))
	var str = make([]byte, 0, len(header)+len(bSlcConLen)+4+conLen)

	str = append(str, header...)
	str = append(str, bSlcConLen...)
	str = append(str, []byte("\r\n\r\n")...)
	for i := strBegin; i < strEnd; i++ {
		str = append(str, req[i])
	}
	conn.Write(str)
}

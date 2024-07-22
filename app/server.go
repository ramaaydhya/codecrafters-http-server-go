package main

import (
	"bytes"
	"fmt"
	"strconv"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func getPath(buff *[]byte) []byte {
	var after []byte
	if isGet(buff) {
		after, _ = bytes.CutPrefix(*buff, []byte("GET "))
	}
	var pathEnd = bytes.IndexByte(after, byte(' '))
	var path = make([]byte, pathEnd)
	copy(path, after[:pathEnd])
	return path
}

func isGet(buff *[]byte) bool {
	return bytes.HasPrefix(*buff, []byte("GET"))
}

func respondGet(conn *net.Conn, status int, respBody *[]byte) (err error) {
	var bodyLen = len(*respBody)
	var conLen = []byte(strconv.Itoa(bodyLen))
	var respLen = 0

	var statusLine = []byte("HTTP/1.1 ")
	switch status {
	case 200:
		statusLine = append(statusLine, []byte("200 OK\r\n")...)
	case 404:
		statusLine = append(statusLine, []byte("404 Not Found\r\n")...)
	}
	respLen += len(statusLine)

	var respHeader = []byte("Content-Type: text/plain\r\nContent-Length: ")
	respHeader = append(respHeader, conLen...)
	respHeader = append(respHeader, []byte("\r\n\r\n")...)
	respLen += len(respHeader)

	respLen += bodyLen

	var resp = make([]byte, 0, respLen)
	resp = append(resp, statusLine...)
	resp = append(resp, respHeader...)
	resp = append(resp, (*respBody)...)
	_, err = (*conn).Write(resp)
	return
}

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
	var path = getPath(&req)
	if isGet(&req) && bytes.Equal(path, []byte("/")) {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}
	if isGet(&req) && bytes.Equal(path, []byte("/user-agent")) {
		var idx = bytes.Index(req, []byte("User-Agent: "))
		var uaBegin = idx + len("User-Agent: ")
		var uaEnd = uaBegin
		for req[uaEnd] != byte('\r') {
			uaEnd++
		}
		var userAgent = make([]byte, uaEnd-uaBegin)
		copy(userAgent, req[uaBegin:uaEnd])
		respondGet(&conn, 200, &userAgent)
		return
	}
	if isGet(&req) && bytes.HasPrefix(path, []byte("/echo")) {
		after, _ := bytes.CutPrefix(path, []byte("/echo/"))
		var body = make([]byte, 0, len(after))
		body = append(body, after...)
		respondGet(&conn, 200, &body)
		return
	}
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

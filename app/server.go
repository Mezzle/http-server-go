package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Println("Starting server on port 4221")
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Panicln("Failed to bind to port 4221")
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Panicln("Failed to close listener")
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln("Error accepting connection: ", err.Error())
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)

	read, err := conn.Read(buffer)
	if err != nil {
		log.Panicln("Error reading: ", err.Error())
	}

	buffer = buffer[:read]
	log.Println("Received: ", string(buffer))

	writeResponse(conn, http.StatusOK, []string{"Content-Length: 5", "Content-Type: text/plain"}, "Hello")
}

func writeLine(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data + "\r\n"))
	if err != nil {
		log.Panicln("Error writing: ", err.Error())
	}
}

func writeHeaders(conn net.Conn, headers []string) {
	for _, header := range headers {
		writeLine(conn, header)
	}
}

func getHttpResponseLine(status int) string {
	return fmt.Sprintf("HTTP/1.1 %d %s", status, http.StatusText(status))
}

func writeStatusLine(conn net.Conn, status int) {
	writeLine(conn, getHttpResponseLine(status))
}

func writeResponse(conn net.Conn, status int, headers []string, body string) {
	writeStatusLine(conn, status)
	writeHeaders(conn, headers)
	writeLine(conn, "")
	writeLine(conn, body)
}

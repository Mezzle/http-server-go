package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
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

	s := string(buffer)

	parts := strings.Split(s, "\n")
	requestLine := parts[0]
	requestLineParts := strings.Split(requestLine, " ")
	path := requestLineParts[1]

	log.Printf("Path: %s", path)

	if path == "/" {
		writeResponse(conn, http.StatusOK, []string{"Content-Type: text/plain"}, "Hello")
		return
	}

	if strings.HasPrefix(path, "/echo/") {
		writeResponse(conn, http.StatusOK, []string{"Content-Type: text/plain"}, path[6:])
		return
	}

	writeResponse(conn, http.StatusNotFound, []string{}, "")
}

func writeLine(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data + "\r\n"))
	if err != nil {
		log.Panicln("Error writing: ", err.Error())
	}
}

func writeHeaders(conn net.Conn, headers []string, contentLength int) {
	for _, header := range headers {
		writeLine(conn, header)
	}
	writeLine(conn, fmt.Sprintf("Content-Length: %d", contentLength))
}

func getHttpResponseLine(status int) string {
	return fmt.Sprintf("HTTP/1.1 %d %s", status, http.StatusText(status))
}

func writeStatusLine(conn net.Conn, status int) {
	writeLine(conn, getHttpResponseLine(status))
}

func writeResponse(conn net.Conn, status int, headers []string, body string) {
	writeStatusLine(conn, status)
	writeHeaders(conn, headers, len(body))
	writeLine(conn, "")
	writeLine(conn, body)
}

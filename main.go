package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Println("Unexpected error")
			break
		}

		if err == io.EOF && len(s) == 0 {
			log.Println("Client disconnected")
			break
		}

		log.Printf("Read %d bytes: %s", len(s), s)

		log.Println("Writing data")
		if _, err := writer.WriteString(s); err != nil {
			log.Fatalln("Unable to write data")
		}
		writer.Flush()

	}
}

func main() {
	// Bind to TCP port 20080 on all interfaces.
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on 0.0.0.0:20080")
	for {
		// Wait for connection. Create net.Conn on connection established.
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		// Handle the connection. Using goroutine for concurrency.
		go echo(conn)
	}
}

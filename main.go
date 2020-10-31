package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err != nil && err != io.EOF {
			log.Println("Unexpected error")
			break
		}

		if err == io.EOF && size == 0 {
			log.Println("Client disconnected")
			break
		}

		log.Printf("Received %d bytes: %s", size, string(b[0:size-2]))

		log.Println("Writing data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func echobuf(conn net.Conn) {
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
		if true {
			go echo(conn)
		} else {
			go echobuf(conn)
		}
	}
}

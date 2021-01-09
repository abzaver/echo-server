package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	useEchoBuf       bool
	ListeningAddress string
	ListeningPort    int
)

func init() {
	flag.BoolVar(&useEchoBuf, "use-buf", false, "using buffering IO")
	flag.StringVar(&ListeningAddress, "listen-address", "0.0.0.0", "listening address")
	flag.IntVar(&ListeningPort, "listen-port", 10500, "listening port")
}

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err != nil && err != io.EOF {
			log.Printf("Unexpected error %s\n", err.Error())
			break
		}

		if err == io.EOF && size == 0 {
			log.Printf("Client %s disconnected\n", conn.RemoteAddr().String())
			break
		}

		log.Printf("Received %d bytes %s", size, string(b))

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
			log.Printf("Unexpected error %s\n", err.Error())
			break
		}

		if err == io.EOF && len(s) == 0 {
			log.Printf("Client %s disconnected\n", conn.RemoteAddr().String())
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
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ListeningAddress, ListeningPort))
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Printf("Listening on %s:%d\n", ListeningAddress, ListeningPort)
	if useEchoBuf {
		log.Println("Using buffering IO...")
	} else {
		log.Println("Using unbuffering IO...")
	}
	for {
		conn, err := listener.Accept()
		log.Printf("Recieved connection from: %s\n", conn.RemoteAddr().String())
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		if useEchoBuf {
			echobuf(conn)
		} else {
			go echo(conn)
		}
	}
}

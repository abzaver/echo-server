package main

import (
	"bufio"
	"net"
	"testing"
)

func TestEcho(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go echo(server)

	message := "Hello, world!\n"
	client.Write([]byte(message))

	response := make([]byte, len(message))
	_, err := client.Read(response)
	if err != nil {
		t.Fatalf("Expected to read '%s' but got error: %v", message, err)
	}

	if string(response) != message {
		t.Fatalf("Expected '%s' but got '%s'", message, string(response))
	}
}

func TestEchoBuf(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go echobuf(server)

	message := "Hello, buffered world!\n"
	client.Write([]byte(message))

	reader := bufio.NewReader(client)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Expected to read '%s' but got error: %v", message, err)
	}

	if response != message {
		t.Fatalf("Expected '%s' but got '%s'", message, response)
	}
}

// Subject
package main

import "fmt"

type server interface {
	handleRequest(string, string) (int, string)
}

// Функция processRequest
func processRequest(s server, url, method string) {
	httpCode, body := s.handleRequest(url, method)
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", url, httpCode, body)
}

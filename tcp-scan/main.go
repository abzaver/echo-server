package main

import (
	"fmt"
	"net"
	"sort"
)

const (
	ports_count   = 16535
	workers_count = 1800
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		//fmt.Printf("%d open\n", p)
		results <- p
	}
}

func main() {
	ports := make(chan int, workers_count)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 0; i < ports_count; i++ {
			ports <- i
		}
	}()

	for i := 0; i < ports_count; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

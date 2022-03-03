package main

import (
	"fmt"
	"sync"
	"time"
)

type empty struct{}

const (
	DATA_SIZE       = 4
	SEMAPHORE_LIMIT = 3
)

func calculate(val int) int {
	fmt.Printf("[%s] Calc for %d\n", time.Now().Format("00:00:00"), val)
	time.Sleep(time.Millisecond * 1200)
	return val * 2
}

func main() {
	data := make([]int, 0, DATA_SIZE)
	for i := 0; i < DATA_SIZE; i++ {
		data = append(data, i+10)
	}
	results := make([]int, DATA_SIZE)
	semaphore := make(chan empty, SEMAPHORE_LIMIT)
	wg := sync.WaitGroup{}

	fmt.Printf("Before: %v\n", data)
	start := time.Now()

	for i, xi := range data {
		wg.Add(1)
		go func(i int, xi int) {
			defer wg.Done()
			semaphore <- empty{}
			results[i] = calculate(xi)
			<-semaphore
		}(i, xi)
	}
	wg.Wait()

	fmt.Printf("After: %v\n", results)
	fmt.Printf("Elapsed: %s\n", time.Since(start))
}

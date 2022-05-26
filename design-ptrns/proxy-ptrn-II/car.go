package main

import "fmt"

type Car struct{}

func (c *Car) Drive() {
	fmt.Println("Car is being driven")
}

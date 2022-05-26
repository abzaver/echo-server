package main

import "fmt"

type CarProxy struct {
	car    Driven
	driver *Driver
}

func NewCarProxy(driver *Driver) Driven {
	return &CarProxy{&Car{}, driver}
}

func (c *CarProxy) Drive() {
	if c.driver.Age >= 16 {
		c.car.Drive()
	} else {
		fmt.Println("Driver too young!")
	}
}

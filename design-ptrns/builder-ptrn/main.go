/*
The Builder pattern is used when the desired product is complex and requires
multiple steps to complete. In this case, several construction methods would be
simpler than a single monstrous constructor. The potential problem with the
multistage building process is that a partially built and unstable product may
be exposed to the client. The Builder pattern keeps the product private until
it’s fully built.

In the below code, we see different types of houses (igloo and normalHouse)
being constructed by iglooBuilder and normalBuilder. Each house type has the
same construction steps. The optional director struct helps to organize the
building process.
*/

// Client code
package main

import "fmt"

func main() {
	normalBuilder := getBuilder("normal")
	iglooBuilder := getBuilder("igloo")

	director := newDirector(normalBuilder)
	normalHouse := director.buildHouse()

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.doorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.windowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.floor)

	director.setBuilder(iglooBuilder)
	iglooHouse := director.buildHouse()

	fmt.Printf("\nIgloo House Door Type: %s\n", iglooHouse.doorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.windowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.floor)

}

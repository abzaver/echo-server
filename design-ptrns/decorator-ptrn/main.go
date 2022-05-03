/*
Decorator is a structural pattern that allows adding new behaviors to objects
dynamically by placing them inside special wrapper objects, called decorators.

Using decorators you can wrap objects countless number of times since both
target objects and decorators follow the same interface. The resulting object
will get a stacking behavior of all wrappers.
*/

// Client code
package main

import "fmt"

func main() {

	pizza := &PizzaMargarita{}

	//Add cheese topping
	pizzaWithCheese := &CheeseTopping{
		pizza: pizza,
	}

	//Add tomato topping
	pizzaWithCheeseAndTomato := &TomatoTopping{
		pizza: pizzaWithCheese,
	}

	fmt.Printf("Price of pizza Margarita with tomato and cheese topping is %d\n", pizzaWithCheeseAndTomato.getPrice())
}

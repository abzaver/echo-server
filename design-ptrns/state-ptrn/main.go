/*
Let’s apply the State Design Pattern in the context of vending machines. For
simplicity, let’s assume that vending machine only has one type of item or
product. Also for simplicity lets assume that a vending machine can be in 4
different states:

* hasItem
* noItem
* itemRequested
* hasMoney

A vending machine will also have different actions. Again for simplicity lets
assume that there are only four actions:

* Select the item
* Add the item
* Insert money
* Dispense item

The State design pattern should be used when the object can be in many
different states and depending upon incoming request the object needs to change
its current state.

In our example, a vending machine can be in many different states and these
states will constantly switch from one to another. Let’s say vending machine
is in itemRequested. Once the action “Insert Money” occurs, the machine
moves to hasMoney state.

Depending on its current state, the machine can behave differently to the same
requests. For example, if user wants to purchase an item then the machine will
proceed if it’s in hasItemState or it will reject if it’s in noItemState.

The code of the vending machine is not polluted with this logics, all the
state-dependent code lives in respective state implementations.
*/

package main

import (
	"fmt"
	"log"
)

func main() {
	vendingMachine := newVendingMachine(1, 10)

	err := vendingMachine.requestItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.addItem(2)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.requestItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

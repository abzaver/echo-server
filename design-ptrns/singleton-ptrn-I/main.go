/*
Singleton is a creational design pattern, which ensures that only one object of
its kind exists and provides a single point of access to it for any other code.

Singleton has almost the same pros and cons as global variables. Although
they’re super-handy, they break the modularity of your code.

You can’t just use a class that depends on a Singleton in some other context,
without carrying over the Singleton to the other context. Most of the time,
this limitation comes up during the creation of unit tests.

Conceptual Example
Usually, a singleton instance is created when the struct is first initialized.
To make this happen, we define the getInstance method on the struct. This
method will be responsible for creating and returning the singleton instance.
Once created, the same singleton instance will be returned every time the
getInstance is called.

How about goroutines? The singleton struct must return the same instance
whenever multiple goroutines are trying to access that instance. Because of
this, it’s very easy to get the singleton design pattern implemented wrong.
The example below illustrates the right way to create a singleton.

Some points worth noting:

There is a nil-check at the start for making sure singleInstance is empty first
time around. This is to prevent expensive lock operations every time the
getinstance method is called. If this check fails, then it means that the
singleInstance field is already populated.

The singleInstance struct is created within the lock.

There is another nil-check after the lock is acquired. This is to ensure that
if more than one goroutine bypasses the first check, only one goroutine can
create the singleton instance. Otherwise, all goroutines will create their own
instances of the singleton struct.
*/

// Client code
package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go getInstance(&wg)
	}

	wg.Wait()
}

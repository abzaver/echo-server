/*
Singleton is a creational design pattern, which ensures that only one object of
its kind exists and provides a single point of access to it for any other code.

Singleton has almost the same pros and cons as global variables. Although
they’re super-handy, they break the modularity of your code.

You can’t just use a class that depends on a Singleton in some other context,
without carrying over the Singleton to the other context. Most of the time,
this limitation comes up during the creation of unit tests.

Another Example
There are other methods of creating a singleton instance in Go:

1. init function
We can create a single instance inside the init function. This is only
applicable if the early initialization of the instance is ok. The init function
is only called once per file in a package, so we can be sure that only a single
instance will be created.

2. sync.Once
The sync.Once will only perform the operation once. See the code in sin
*/

// Client code
package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go getInstance(&wg)
	}

	wg.Wait()
}

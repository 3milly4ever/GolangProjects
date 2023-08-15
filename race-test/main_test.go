package main

import (
	"sync/atomic"
	"testing"
)

// a race condition is when two threads or routines are trying to read or modify the same variable
// do a race test to make sure you dont read or write to the same variable from different routines
func TestDataRaceConditions(t *testing.T) {
	//we usually have some state, we do some api queries and goroutines, we fetch some data and we do some logic with the state
	//a map is not concurrent safe, so if there is a map inside a goroutine there will be a race condition 100%
	var state int32

	//state 2
	//state 3
	//state 4
	//if you have multiple states, mutexes can cause problems
	//if you have a lock for state 1, the other routine is willing to write to state2 but it's the same lock.
	//you technically wouldn't be able to have a race condition because state 1 and state 2 are different structures,
	//but if you have the same lock for different states you run into a deadlock. mutexes will slow your program down too.
	//you can use unbuffered channels to communicate so you dont need a mutex.

	//if you are using mutexes, call them either mu or lock

	//var mu sync.RWMutex //read write mutex, a mutex can lock goroutines, the goroutine that holds the lock

	//can read and write depending on the lock.
	//the mutex that holds the lock to a certain goroutine can read and write to the variable, if someone else without a mutex
	//tries to access the variable in the goroutine needs to wait because he doesn't have the lock

	//atomic value makes sure that the value can only be adjusted by one
	for i := 0; i < 10; i++ {
		go func(i int) {
			//state += int32(i)
			//the atmoic addint is the same as the above code, but just atomically.
			atomic.AddInt32(&state, int32(i))
			//mu.Lock()
			//you can do mu.RLock() which is a read lock

			//business logic can read or write to the variable
			//mu.Unlock() //once it is unlocked now the mu.Lock can be used elsewhere
		}(i)

	}

}

func main() {

}

package c1

import (
	"fmt"
	"sync"
)

type Chapter1 struct{}

func (c *Chapter1) RunChapter() {
	// raceCondition()
	// syncMemoryAccess()
	// atomicity()
	// deadLock()
	// liveLock()
	starvation()
}

// Race condition - when two or more threads try to modify a shared resource at the same time.
func raceCondition() {
	// critical section - the part of the program that accesses shared resources.
	var data int
	go func() {
		data++ // critical section 1 - incrementing the data
	}()
	if data == 0 { // critical section 2 - reading the data
		fmt.Printf("the value is %v.\n", data) // critical section 3 - the value is 0.
	}
}

// This solves our data race, but not our race condition, since whether the value is printed as 0 or 1 is still dependent on the order of execution
func syncMemoryAccess() {
	var memoryAccess sync.Mutex
	var value int
	go func() {
		memoryAccess.Lock()
		value++
		memoryAccess.Unlock()
	}()
	memoryAccess.Lock()
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock()
}

// Atomicity - the property of an operation to be carried out as a single unit of execution without any interference by another thread.
// increment is not atomic - its actually 3 operations: read, increment, write.
func atomicity() {
	i := 0
	i++
	fmt.Println(i)
}

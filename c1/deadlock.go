package c1

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

// deadlock is defined by the 4 conditions:
// 1. Mutual exclusion - only one thread can access a resource at a time.
// 2. Hold and wait - a thread holds a resource and waits for another.
// 3. No preemption - a resource can only be released by the thread holding it.
// 4. Circular wait - a cycle of threads waiting for resources.
func deadLock() {
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()
		println("sum is", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	// this gives deadlock because:
	go printSum(&a, &b) // a is locked, then in the 2 seconds it has to wait before acquiring b, b is already locked by the other goroutine
	go printSum(&b, &a) // technically there is a race condition here because its possible that the second goroutine doesn't start until after 2 seconds
	wg.Wait()
}

func liveLock() {
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(dir, 1)
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
		}
		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	// this code tries to get people to move at the same cadence (same time), and if they are unable to, they switch directions and try again.
	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, " %v gave up.", name)
	}

	// what happens here though is if you have 2 people in the hallway, they try to both go left, find that theres another person there,
	// and then try to go right, find that theres another person there, and then try to go left again, and so on.
	// Hence, a live lock.
	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Bob")
	peopleInHallway.Wait()
}

// livelock is subset of starvation (when a concurrent process cannot get all the resources it needs)
// livelock is when all concurrent processes are starved equally, and NO work is accomplished
// however, there are other cases where some work is accomplished, but not all is accomplished as efficiently as possible
func starvation() {
	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()
		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}
		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()
		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}

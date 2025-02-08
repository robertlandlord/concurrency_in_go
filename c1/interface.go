package c1

// when writing functions that are concurrent, comments need to address the following:
// 1. who is responsible for the concurrency
// 2. how is the problem space mapped onto concurrency primitives?
// 3. who is responsible for the synchronization?
// e.g.
// CalculatePi calculates digits of pi between the begin and end place.
// Intenrally, CalculatePi will create FLOOR((end-begin)/2) concurrent
// processes which recursively call CalculatePi. Synchronization of
// writes to pi are handled internally by the Pi struct
func CalculatePi(begin, end int64, pi *Pi) {}

type Pi struct{}

// when modeling the signature of the function, we can also give signals to what is happening
// e.g. for concurrency + synchronization, returning a read-only channel suggests that
// both concurrency and synchronization is handled internally by the function.
func CalculatePiBetter(begin, end int64) <-chan uint {
	result := make(chan uint)

	go func() {
		defer close(result)
		for i := begin; i < end; i++ {
			// Mock calculation, just sending the index as a result
			result <- uint(i)
		}
	}()

	return result
}

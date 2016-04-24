package bevel

import "sync/atomic"

// Counter is used to count the number of messages
// that the Listener has processed since it was started.
type Counter struct {
	msgCounter uint64
}

// Get is a getter for the value of msgCounter.
func (c *Counter) Get() uint64 {
	return c.msgCounter
}

// Inc is a thread-safe incrementer of the value of msgCounter.
func (c *Counter) Inc() {
	atomic.AddUint64(&c.msgCounter, 1)
}

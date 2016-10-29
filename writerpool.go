package bevel

import (
	"fmt"
	"log"
	"sync"
)

// Writer is an interface that defines the operations
// that Writer implementors must adhere to.
// This is used by WriterPool.Write() to write the message contents.
type Writer interface {
	Write(Message) error
}

// WriteCloser is an interface that combines the operations
// of Writer and Closer.
type WriteCloser interface {
	Writer
	Closer
}

// WriterPool is a thread-safe pool of Writer's.
type WriterPool struct {
	sync.Mutex
	writers []Writer // interfaces and pointers: http://stackoverflow.com/questions/13511203/why-cant-i-assign-a-struct-to-an-interface
}

// AddWriter adds a Writer to the pool of Writer's.
func (bewp *WriterPool) AddWriter(w Writer) {
	bewp.Lock()
	defer bewp.Unlock()

	bewp.writers = append(bewp.writers, w)
}

// Writer iterates through the pool of Writer's and
// and calls Write() on each one of them in a thread-safe manner.
// Thread-safety is achieved independently per Writer.
func (bewp *WriterPool) Write(m Message) {
	var wg sync.WaitGroup

	// nice articles:
	// - http://blog.launchdarkly.com/golang-pearl-thread-safe-writes-and-double-checked-locking-in-go/
	// - https://golang.org/pkg/sync/#example_WaitGroup
	for _, w := range bewp.writers {
		wg.Add(1)
		go func(w Writer, m Message) {
			defer wg.Done()
			if err := w.Write(m); err != nil {
				log.Printf("ERROR - failed to write message: %s\n", err) // TODO add details of writer and message
			}
		}(w, m)
	}

	wg.Wait()
}

// String is an implementation of Stringer.
func (bewp *WriterPool) String() string {
	bewp.Lock()
	defer bewp.Unlock()

	s := "Registered writers:"
	for _, w := range bewp.writers {
		s += fmt.Sprintf(" %T", w)
	}

	return s
}

// NewWriterPool creates a new WriterPool.
func NewWriterPool() *WriterPool {
	wp := WriterPool{writers: []Writer{}}

	return &wp
}

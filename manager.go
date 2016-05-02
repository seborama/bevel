package bevel

import (
	"fmt"
	"time"
)

// Manager holds the properties needed by the business event logger.
// It is the receiver of utility methods for consumers.
type Manager struct {
	done        chan bool
	bus         chan MessageEnvelop
	writersPool WriterPool
	msgCounter  Counter
}

// StartNewListener creates a new business event bus and adds
// a Writer to the WriterPool.
func StartNewListener(w Writer) *Manager {
	wp := *NewWriterPool()
	wp.AddWriter(w)

	bem := Manager{
		make(chan bool),
		make(chan MessageEnvelop),
		wp,
		Counter{0},
	}

	go bem.listen()

	return &bem
}

// Done sends the listener a message to terminate.
func (bem *Manager) Done() {
	if bem.done == nil {
		return
	}

	bem.done <- true // Sending listen() goroutine "termination Ping"
	<-bem.done       // Waiting for listen() goroutine to respond with "termination Pong"

	close(bem.bus)
	bem.bus = nil

	close(bem.done)
	bem.done = nil
}

// MessageEnvelop wraps Message with additional info.
type MessageEnvelop struct {
	ProcessedTSUnixNano int64
	Message             Message
}

// Post sends a Message to the business event message bus
// for ingestion by all Writer's in the WriterPool.
func (bem *Manager) Post(m Message) {
	if bem.bus == nil || bem.done == nil {
		return
	}

	me := MessageEnvelop{
		ProcessedTSUnixNano: time.Now().UnixNano(),
		Message:             m,
	}

	bem.bus <- me
}

// AddWriter adds a Writer to the WriterPool.
func (bem *Manager) AddWriter(w Writer) {
	bem.writersPool.AddWriter(w)
}

func (bem *Manager) writeMessage(m Message) {
	bem.writersPool.Write(m)
}

// String implements Stringer.
func (bem *Manager) String() string {
	s := bem.writersPool.String()
	s += fmt.Sprintf(" - Total number of messages posted: %d", bem.msgCounter.Get())

	return fmt.Sprintf("%s", s)
}

// listen is the main loop of the business event loop.
func (bem *Manager) listen() {
	defer func() {
		bem.done <- true // Sending "Termination Pong" response
	}()

	loop := true
	for loop {
		select {
		case m := <-bem.bus:
			// Received a Message.
			bem.msgCounter.Inc()
			bem.writeMessage(m)
		case <-bem.done:
			// Received "Termination Ping" request.
			loop = false
			break
		}
	}
}

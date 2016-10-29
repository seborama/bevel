package bevel

import (
	"errors"
	"fmt"
	"time"
)

// Closer is an interface that defines the operations
// that Closer implementors must adhere to.
type Closer interface {
	Close() error
}

// Poster is an interface that defines the operations
// that Poster implementors must adhere to.
type Poster interface {
	Post(Message)
}

// AddWriterer is an interface that defines the operation
// AddWriter.
type AddWriterer interface {
	AddWriter(Writer)
}

// EventBusManager is an interface that combines the operations
// of an Event Bus Manager.
type EventBusManager interface {
	AddWriterer
	Poster
	Closer
	fmt.Stringer
}

// Manager holds the properties needed by the business event logger.
// It is the receiver of utility methods for consumers.
type Manager struct {
	done        chan bool
	bus         chan MessageEnvelop
	writersPool WriterPool
	msgCounter  Counter
}

// MessageEnvelop wraps Message with additional info.
type MessageEnvelop struct {
	ProcessedTSUnixNano int64
	Message             Message
}

// StartNewListener creates a new business event bus and adds
// a Writer to the WriterPool.
func StartNewListener(w WriteCloser) EventBusManager {
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

// Post sends a Message to the business event message bus
// for ingestion by all Writer's in the WriterPool.
func (bem *Manager) Post(m Message) {
	// ensure the bus is open for messages (i.e. "post office is open")
	if bem.bus == nil || bem.done == nil {
		return
	}

	// wrap the application message into an envelop (i.e. "put the letter in an envelop and affix a stamp")
	me := MessageEnvelop{
		ProcessedTSUnixNano: time.Now().UnixNano(),
		Message:             m,
	}

	// post the envelop on the bus (i.e. "post the letter")
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

ListenerLoop:
	for {
		select {
		case m := <-bem.bus:
			// Received a Message wrapped in a MessageEnvelop.
			// Call the writer to write it to destination - in this case to Kafka.
			bem.msgCounter.Inc()
			bem.writeMessage(m)
		case <-bem.done:
			// Received "Termination Ping" request.
			break ListenerLoop
		}
	}
}

// Close closes the channels in the Manager.
func (bem *Manager) Close() error {
	if bem.done == nil {
		return errors.New("this event bus manager is already closed")
	}

	bem.done <- true // Sending listen() goroutine "termination Ping"
	<-bem.done       // Waiting for listen() goroutine to respond with "termination Pong"

	close(bem.bus)
	bem.bus = nil

	close(bem.done)
	bem.done = nil

	return nil
}

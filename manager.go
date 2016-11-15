package bevel

import (
	"errors"
	"fmt"
	"log"
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
	bus         chan Message
	writersPool WriterPool
	msgCounter  Counter
}

// StartNewListener creates a new business event bus and adds
// a Writer to the WriterPool.
func StartNewListener(w WriteCloser) EventBusManager {
	wp := *NewWriterPool()
	wp.AddWriter(w)

	bem := Manager{
		done:        make(chan bool),
		bus:         make(chan Message),
		writersPool: wp,
		msgCounter:  Counter{0},
	}

	go bem.listen()

	return &bem
}

// Post sends a Message to the business event message bus
// for ingestion by all Writer's in the WriterPool.
func (bem *Manager) Post(m Message) {
	// ensure the bus is open for messages (i.e. "post office is open")
	if bem.bus == nil || bem.done == nil {
		log.Printf("the event bus is closed - lost message: %#v", m)
		return
	}

	bem.bus <- NewMesageEnvelop(m)
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
	s += fmt.Sprintf(" - total number of messages posted: %d", bem.msgCounter.Get())

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
		case m, ok := <-bem.bus:
			if ok {
				// Received a Message wrapped in a MessageEnvelop.
				// Call the writer to write it to destination - in this case to Kafka.
				bem.msgCounter.Inc()
				bem.writeMessage(m)
			}
		case <-bem.done:
			// Received "Termination Ping" request.
			// Drain the remaining messages on the bus and break out.
			for m := range bem.bus {
				bem.msgCounter.Inc()
				bem.writeMessage(m)
			}
			break ListenerLoop
		}
	}
}

// Close closes the channels in the Manager.
// The recommended approach is for a channel to be used unidirectionally and
// be closed by the sender rather than the receivers.
// This means that it is the responsibility of the Posters to close the
// event bus when no more messages are being posted.
// See an example implementation in main_test using a sync.WaitGroup.
func (bem *Manager) Close() error {
	if bem.done == nil {
		return errors.New("this event bus manager is already closed")
	}

	close(bem.bus)
	bem.done <- true
	<-bem.done
	bem.bus = nil

	close(bem.done)
	bem.done = nil

	return nil
}

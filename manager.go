package bevel

import (
	"errors"
	"fmt"
	"log"
	"sync"
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
	wg          sync.WaitGroup
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
		wg:          sync.WaitGroup{},
	}

	go bem.listen()

	return &bem
}

// Post sends a Message to the business event message bus
// for ingestion by all Writer's in the WriterPool.
func (bem *Manager) Post(m Message) {
	defer func() {
		if r := recover(); r != nil {
			// manual attempt to send write the message
			log.Println("error posting to the event bus - sending synchronously:", r)
			bem.msgCounter.Inc()
			bem.writeMessage(NewMesageEnvelop(m))
		}

		bem.wg.Done()
	}()

	bem.wg.Add(1)

	// post the envelop on the bus (i.e. "post the letter")
	if bem.bus == nil || bem.done == nil {
		panic("the event bus is closed")
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
	defer func() { bem.done <- true }()

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
			break ListenerLoop
		}
	}
}

// Close closes the channels in the Manager.
func (bem *Manager) Close() error {
	if bem.done == nil {
		return errors.New("this event bus manager is already closed")
	}

	// wait for all Posters to complete
	bem.wg.Wait()

	close(bem.bus)
	bem.done <- true
	<-bem.done
	bem.bus = nil

	close(bem.done)
	bem.done = nil

	// pick up the Posters that may have spawned before we got a chance to close the channels.
	bem.wg.Wait()

	return nil
}

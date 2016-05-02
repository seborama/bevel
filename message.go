package bevel

// Message is a high level interface that all
// Business Events must implement.
type Message interface{}

// StandardMessage is the first implementor of Message.
// This represents acts as a header for all Message's.
type StandardMessage struct {
	EventName         string
	CreatedTSUnixNano int64
}

package bevel

import "time"

// Message is a high level interface that all
// Business Events must implement.
type Message interface{}

// StandardMessage is the first implementor of Message.
// This represents acts as a header for all Message's.
type StandardMessage struct {
	EventName         string
	CreatedTSUnixNano int64
}

// MessageEnvelop wraps Message with additional info.
type MessageEnvelop struct {
	Message
	ProcessedTSUnixNano int64
}

// NewMesageEnvelop creates a new MessageEnvelop from a supplied message.
func NewMesageEnvelop(m Message) *MessageEnvelop {
	return &MessageEnvelop{
		Message:             m,
		ProcessedTSUnixNano: time.Now().UnixNano(),
	}
}

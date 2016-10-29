package bevel

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrUnableToTranscodeMessage occurs when the content of the message
// could not be transcoded by Write() for output.
var ErrUnableToTranscodeMessage = errors.New("unable to transcode message")

// ConsoleBEWriter is a simple Writer implementation..
type ConsoleBEWriter struct{}

// Write outputs the contents of Message to the console.
func (bew *ConsoleBEWriter) Write(m Message) error {
	// TODO - change the hard-coding of json.Marshal() to a
	// strategy pattern.
	// This would allow to inject a strategy to write messages.
	// For instance, a JSONWriterStrategy or a MapWriterStrategy (which
	// would return a map) or an XMLWriterStrategy, etc.
	json, err := json.Marshal(m)
	if err != nil {
		return ErrUnableToTranscodeMessage
	}

	_, err = fmt.Printf("%s\n", string(json))
	return err
}

// Close does not perform any action in the case of ConsoleBEWriter.
func (bew *ConsoleBEWriter) Close() error {
	return nil
}

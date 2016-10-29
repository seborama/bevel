package bevel_test

import (
	"os"
	"testing"
	"time"

	"github.com/seborama/bevel"
)

// CounterMsg is an example of message that can be sent
// on the Listener's bus.
type CounterMsg struct {
	bevel.StandardMessage
	Counter int
}

func TestBevel(t *testing.T) {
	// Start new BEManager with a Console Business Event Writer.
	bem := bevel.StartNewListener(&bevel.ConsoleBEWriter{})
	defer func() {
		bem.Close()
	}()

	// Add another Console Business Event Writer just to show it can be done.
	bem.AddWriter(&bevel.ConsoleBEWriter{})

	// Create some business events
	for i := 1; i <= 5; i++ {
		m := CounterMsg{
			StandardMessage: bevel.StandardMessage{
				EventName:         "test_event",
				CreatedTSUnixNano: time.Now().UnixNano(),
			},
			Counter: i,
		}

		bem.Post(m)
	}
}

func TestMain(m *testing.M) {
	v := m.Run()
	os.Exit(v)
}

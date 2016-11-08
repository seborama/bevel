package bevel_test

import (
	"os"
	"testing"
	"time"

	"log"

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
	j := 0
	for i := 1; i <= 5; i++ {
		j = i
		m := CounterMsg{
			StandardMessage: bevel.StandardMessage{
				EventName:         "test_event",
				CreatedTSUnixNano: time.Now().UnixNano(),
			},
			Counter: i,
		}

		go bem.Post(m)
	}

	time.Sleep(time.Millisecond * 500)
	s := bem.String()
	if s != "Registered writers: *bevel.ConsoleBEWriter *bevel.ConsoleBEWriter - total number of messages posted: 5" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}
	bem.Close()

	for i := j + 1; i <= j+5; i++ {
		m := CounterMsg{
			StandardMessage: bevel.StandardMessage{
				EventName:         "test_event",
				CreatedTSUnixNano: time.Now().UnixNano(),
			},
			Counter: i,
		}

		go bem.Post(m)
	}

	time.Sleep(time.Millisecond * 500)
	log.Println(bem)
	s = bem.String()
	if s != "Registered writers: *bevel.ConsoleBEWriter *bevel.ConsoleBEWriter - total number of messages posted: 10" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}
}

func TestMain(m *testing.M) {
	v := m.Run()
	os.Exit(v)
}

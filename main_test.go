package bevel_test

import (
	"os"
	"testing"
	"time"

	"log"

	"sync"

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

		log.Println(bem)
		s := bem.String()
		if s != "Registered writers: *bevel.ConsoleBEWriter *bevel.ConsoleBEWriter - total number of messages posted: 5" {
			t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
		}
	}()

	// Add another Console Business Event Writer just to show it can be done.
	bem.AddWriter(&bevel.ConsoleBEWriter{})

	// Create some business events
	wg := sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		m := CounterMsg{
			StandardMessage: bevel.StandardMessage{
				EventName:         "test_event",
				CreatedTSUnixNano: time.Now().UnixNano(),
			},
			Counter: i,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			bem.Post(m)
		}()
	}

	wg.Wait()
}

func TestMain(m *testing.M) {
	v := m.Run()
	os.Exit(v)
}

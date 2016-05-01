package bevel_test

import (
	"testing"

	"github.com/seborama/bevel"
)

type testWriter3 struct{}

func (r testWriter3) Write(m bevel.Message) error {
	return nil
}

func TestWPAddWriter(t *testing.T) {
	w3 := testWriter3{}
	wp := bevel.NewWriterPool()
	wp.AddWriter(w3)

	s := wp.String()
	if s != "Registered writers: bevel_test.testWriter3" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}
}

func TestNewWriterPool(t *testing.T) {
	// Tested in TestWPAddWriter()
}

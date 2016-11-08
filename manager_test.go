package bevel_test

import (
	"testing"

	"github.com/seborama/bevel"
)

type testWriter1 struct{}
type testWriter2 struct{ testWriter1 }

func (r testWriter1) Write(m bevel.Message) error {
	return nil
}

func (r testWriter1) Close() error {
	return nil
}

func TestStartNewListener(t *testing.T) {
	// one writer, no messages
	w := testWriter1{}
	bem := bevel.StartNewListener(w)
	if bem == nil {
		t.Error("StartNewListener should have returned a Manager but it returned nil")
	}

	s := bem.String()
	if s != "Registered writers: bevel_test.testWriter1 - total number of messages posted: 0" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}

	// one writer, two messages
	m := bevel.StandardMessage{}
	bem.Post(m)
	bem.Post(m)
	bem.Close() // ensure all messages have been processed
	s = bem.String()
	if s != "Registered writers: bevel_test.testWriter1 - total number of messages posted: 2" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}

	// one message past call to Done()
	// should not timeout & panic
	// should not accept the message (internal msgCounter unchanged)
	bem.Post(m)
	s = bem.String()
	if s != "Registered writers: bevel_test.testWriter1 - total number of messages posted: 3" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}
}

func TestClose(t *testing.T) {
	// Tested in TestStartNewListener()
}

func TestPost(t *testing.T) {
	// Tested in TestStartNewListener()
}

func TestString(t *testing.T) {
	// Tested in TestStartNewListener() & TestAddWriter()
}

func TestAddWriter(t *testing.T) {
	w1 := testWriter1{}
	bem := bevel.StartNewListener(w1)
	if bem == nil {
		t.Error("StartNewListener should have returned a Manager but it returned nil")
	}
	defer bem.Close()

	w2 := testWriter2{}
	bem.AddWriter(w2)
	s := bem.String()
	if s != "Registered writers: bevel_test.testWriter1 bevel_test.testWriter2 - total number of messages posted: 0" {
		t.Errorf("Manager is not in the expected state after call to StartNewListener: %s\n", s)
	}
}

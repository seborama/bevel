package bevel_test

import (
	"testing"

	"github.com/seborama/bevel"
)

func TestGet(t *testing.T) {
	c := bevel.Counter{}
	v := c.Get()
	if v != 0 {
		t.Errorf("Get should have returned 0 but it returned %d\n", c)
	}
}

func TestInc(t *testing.T) {
	c := bevel.Counter{}
	c.Inc()
	v := c.Get()
	if v != 1 {
		t.Errorf("After one call to Set, Get should have returned 1 but it returned %d\n", c)
	}

	for i := 1; i <= 10; i++ {
		c.Inc()
	}
	v = c.Get()
	if v != 11 {
		t.Errorf("After 10 additionals call to Set, Get should have returned 11 but it returned %d\n", c)
	}
}

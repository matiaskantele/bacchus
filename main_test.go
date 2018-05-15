package main

import "testing"

func TestPlaceholder(t *testing.T) {
	expected := 5
	got := 2+3
	if expected != got {
		t.Errorf("Sum test [%s], expected [%d], actual [%d]", "2+3", expected, got)
	}
}

package main

import "testing"

func TestSum(t *testing.T) {
	expected := 5
	got := sum(2, 3)
	if expected != got {
		t.Errorf("Sum test [%s], expected [%d], actual [%d]", "2+3", expected, got)
	}
}

package src

import "testing"

func TestHandle(t *testing.T) {
	got := Handle("x")
	if got != "handled: x" {
		t.Fatalf("got %q", got)
	}
}

package main

import "testing"

func TestHello(t *testing.T) {
	got := hello()
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

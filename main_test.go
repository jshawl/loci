package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	pause := time.Duration(rand.Int63n(3000)+100) * time.Millisecond // nolint:gosec
	time.Sleep(pause)
	t.Fatalf("yikes!")
}

package main

import (
	"strings"
	"testing"
)

func TestSomething(t *testing.T) {
	t.Fatalf(strings.Repeat("yikes!\n", 100))
}

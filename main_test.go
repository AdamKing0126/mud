package main

import (
	"strings"
	"testing"
)

func TestGetPlayerInput(t *testing.T) {
	reader := strings.NewReader("Test Player\n")
	expected := "Test Player"

	result := getPlayerInput(reader)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

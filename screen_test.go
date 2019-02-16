package main

import "testing"

func TestNewScreen(t *testing.T) {
	screen := NewScreen(10.0, 10.0)
	if screen.Puck.position.AtVec(0) > 10.0 ||
		screen.Puck.position.AtVec(0) < 0 ||
		screen.Puck.position.AtVec(1) < 0 ||
		screen.Puck.position.AtVec(1) > 10 {
		t.Error("Puck out of bounds")
	}
}

func Test
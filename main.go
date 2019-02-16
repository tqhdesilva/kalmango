package main

import "time"

func main() {
	var screen *Screen
	screen = NewScreen(10, 10)
	c := make(chan time.Time)
	go screen.Run(.1, c)
}

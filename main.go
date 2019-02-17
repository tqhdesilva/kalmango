package main

import "time"

func main() {
	screen := NewScreen(10, 10)
	c := make(chan time.Time)
	go screen.Run(.1, c)
}

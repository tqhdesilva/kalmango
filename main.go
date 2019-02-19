package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	screen := NewScreen(10, 10)
	c := make(chan time.Time)
	go screen.Run(.1, c)

}

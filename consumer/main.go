package main

import (
	"time"
)

func main() {
	go Listen()

	for {
		time.Sleep(5 * time.Second)
	}
}

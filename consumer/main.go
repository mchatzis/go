package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello consumer")
		time.Sleep(5 * time.Second)
	}
}

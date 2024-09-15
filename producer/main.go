package main

import (
	"fmt"
	"time"
)

func main() {

	for {
		fmt.Println("Hello producer")
		time.Sleep(5 * time.Second)
	}
}

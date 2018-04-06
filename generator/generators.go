package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string { // Returns recieve-only channel of strings
	c := make(chan string)
	go func() { // We launch the gorotuine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller
}

func main() {
	c := boring("boring!") // Function returning a channel
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You are boring! I am leaving")
}

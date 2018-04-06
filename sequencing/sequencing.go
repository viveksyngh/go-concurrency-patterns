package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

func boring(msg string) <-chan Message { // Returns recieve-only channel of strings
	c := make(chan Message)
	waitForIt := make(chan bool)
	go func() { // We launch the gorotuine from inside the function
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c // Return the channel to the caller
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

func main() {
	c := fanIn(boring("Joe"), boring("Ann"))
	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're all boring! I'm leaving.")

}

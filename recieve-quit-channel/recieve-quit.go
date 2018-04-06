package main

import (
	"fmt"
	"math/rand"
)

func cleanup() {

}

func boring(message string, quit chan string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s : %d", message, i):
				// time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-quit:
				cleanup()
				quit <- "See You!"
				return
			}
		}
	}()
	return c
}

func main() {
	quit := make(chan string)
	c := boring("Joe", quit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	quit <- "Bye!"
	fmt.Println(<-quit)
}

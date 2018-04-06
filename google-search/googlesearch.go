package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

//Result custom type for search result
type Result string

//Search type for the different search type
type Search func(query string) Result

//Google first version of google search
func Google(query string) (results []Result) {
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return results
}

//Google2 version 2.0 of google search which uses goroutines to get search result for each type
func Google2(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return results
}

//Google21 version 2.1 of google search which has timeout of 80 ms for each search type to result
func Google21(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Timed out")
			return
		}
	}
	return results
}

//First returns result of the search from the fastest replica
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

//Google3 Enhanced version of google search which runs multiple replicas for each search type
func Google3(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- First(query, fakeSearch("Web1"), fakeSearch("Web2")) }()
	go func() { c <- First(query, fakeSearch("Image1"), fakeSearch("Image2")) }()
	go func() { c <- First(query, fakeSearch("Video1"), fakeSearch("Video2")) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Timed Out")
			return
		}
	}
	return
}

//fakeSearch it takes a search query and fakes the search with
func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	// results := Google("Golnag")
	// results := Google2("Golang")
	// results := Google21("Golang")
	// results := First("Golang", fakeSearch("replica 1"), fakeSearch("replica 2"))
	results := Google3("Golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

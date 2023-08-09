package main

//make ansynchronous calls and aggregate data
//you might have 10 different api calls to 10 different 3rd-parties and they all take 100 miliseconds, you will wait a long time if it's one-by-one
//if you do it like this, it will take a maximum 100 miliseconds

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	now := time.Now()
	userID := 10

	//communicate with go channel
	//create go response channel that will be a string
	//in a real world scenario, it will be your own structure that holds a certain type of data
	//after receiving all the data from the goroutines it closes, processes the received data and performs any remaining tasks
	respch := make(chan string, 128) //buffered channel because there is capacity

	//a wait group.
	wg := &sync.WaitGroup{}

	//put the channel inside a function
	go fetchUserData(userID, respch, wg)
	go fetchUserRecommendations(userID, respch, wg)
	go fetchUserLikes(userID, respch, wg)

	//until we say done, wait will wait.
	//if no one says done, there is a deadlock
	wg.Add(3) // because we have three go channels, informs the waitgroup there are 3 goroutines we need to wait for.
	//we wait to close the main goroutine (respch)
	wg.Wait() //blocks the main goroutine until all three have completed their work and called wg.Done(). once they do it unblocks.

	//	wg.Done()

	close(respch) //closes the channel meaning it will not receive any more data

	//you can range over the go channel we create above
	for resp := range respch {
		fmt.Println(resp)
	}

	fmt.Println(time.Since(now))
}

func fetchUserData(userID int, respch chan string, wg *sync.WaitGroup) {
	//we are mocking how long it takes to retrieve data from a database
	time.Sleep(80 * time.Millisecond)
	//respond
	respch <- "user data"
	//we are done
	wg.Done()
}

// our 'weakest link' because it takes the longest to fetch/process
// so even if the rest of the go channels finish fast, we would still have to wait for this one to finish making all of them this slow
func fetchUserRecommendations(userID int, respch chan string, wg *sync.WaitGroup) {
	//we are mocking how long it takes to fetch data that AI is involved with
	time.Sleep(120 * time.Millisecond)

	respch <- "user recommendations"

	wg.Done()
}

func fetchUserLikes(userID int, respch chan string, wg *sync.WaitGroup) {
	//faster if it's a low amount of likes
	time.Sleep(50 * time.Millisecond)

	respch <- "user likes"

	wg.Done()
}

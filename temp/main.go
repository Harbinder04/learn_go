package main

import (
	"fmt"
	// "sync"
	"time"
)


type job struct {
	Type string
	payload any
}
func main() {
	jobChan := make(chan job)
	shutsig := make(chan struct{})
	// var wg sync.WaitGroup

	go worker(jobChan, shutsig)

	go func() {
		jobChan <- job{Type: "W", payload: "har"}
		jobChan <- job{Type: "W", payload: "har"}
		jobChan <- job{Type: "W", payload: "har"}
		time.Sleep(1 * time.Second)
		close(shutsig)
	}()

	fmt.Print("I am here")
	time.Sleep(8 * time.Second)
	fmt.Print("close")

}

func worker(jobChan <-chan job, done <-chan struct{}) {
	for {
		select {
		case job := <-jobChan :
			process(job)
		case <-done: 
			return
		}
	}
}

func process(st job) {
	fmt.Print(st)
	time.Sleep(2 * time.Second)
}
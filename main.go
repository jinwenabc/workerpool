package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	payload int
}

var maxJobs = flag.Int("maxjobs", 10, "specified max jobs to be cached")
var maxWorkers = flag.Int("maxworkers", 10, "specified max workers concurrent working")

func enQueue(jobQueue chan Job, reqs ...int)  {
	go func() {
		for _, req := range reqs{
			job := Job{payload:req}
			jobQueue <- job
		}
	}()
}

type Worker struct {
	JobC chan Job
	quit chan struct{}
}

func NewWorker() Worker {
	return Worker{
		JobC: make(chan Job),
		quit: make(chan struct{}),
	}
}

func (w Worker)Start(workerPool chan Worker)  {
	go func() {
		for{
			// Register service
			workerPool <- w
			select {
			case job:= <- w.JobC:
				fmt.Printf("%v, dealt with payload: %d\n", time.Now(), job.payload)
			case <- w.quit:
				return
			}
		}
	}()
}

func (w Worker)Stop()  {
	go func() {
		w.quit <- struct{}{}
	}()
}

func DispatchJobs(workerPool chan Worker, jobQueue chan Job)  {
	for{
		job := <- jobQueue
		go func(job Job) {
			worker := <- workerPool
			worker.JobC <- job
		}(job)
	}
}

func main()  {
	flag.Parse()
	jobQueue := make(chan Job, *maxJobs)
	workerPool := make(chan Worker, *maxWorkers)
	for i:=1; i<=*maxWorkers; i++{
		NewWorker().Start(workerPool)
	}

	go DispatchJobs(workerPool, jobQueue)
	for {
		time.Sleep(time.Second)
		payloads := make([]int, 0)
		for i:=0; i<20; i++{
			payloads = append(payloads, rand.Intn(100))
		}
		enQueue(jobQueue, payloads...)
	}
}
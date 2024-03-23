package job

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	Wg         *sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	queue := make(chan Job)
	return &Dispatcher{
		WorkerPool: pool,
		JobQueue:   queue,
		Wg:         &sync.WaitGroup{},
	}
}

func (d *Dispatcher) AddJob(job Job) {
	d.Wg.Add(1)
	d.JobQueue <- job
}

func (d *Dispatcher) Close() {
	close(d.JobQueue)
}

func (d *Dispatcher) Run() {
	ready := &sync.WaitGroup{}
	ready.Add(cap(d.WorkerPool)) // Ajouter le nombre de Workers au WaitGroup

	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.JobQueue, d.Wg)
		worker.Start(ready)
		d.WorkerPool <- worker.JobQueue
	}

	for {
		select {
		case job, ok := <-d.JobQueue:
			if ok {
				fmt.Println("Dispatcher: received job")
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}
		default:
			// Si la JobQueue est vide, sortir de la boucle
			if len(d.JobQueue) == 0 {
				return
			}
		}
	}
}

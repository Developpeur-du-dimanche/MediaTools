package job

import (
	"sync"
)

type Job struct {
	Name string
	Work func() (bool, error)
}

type Worker struct {
	JobQueue chan Job
	Wg       *sync.WaitGroup
}

func NewWorker(jobQueue chan Job, wg *sync.WaitGroup) *Worker {
	return &Worker{JobQueue: jobQueue, Wg: wg}
}

func (w *Worker) Start(ready *sync.WaitGroup) {
	ready.Done() // Indiquer que le Worker est prêt
	go func() {
		for job := range w.JobQueue {
			job.Work()
			w.Wg.Done()
		}
	}()
}

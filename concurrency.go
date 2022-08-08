package GoParallel

import (
	"errors"
	"sync"
)

var ErrWorkersReturnWithError error = errors.New("some workers return with error")

type Processor struct {
	errLock      *sync.Mutex
	progressLock *sync.Mutex
	progress     int
	wg           sync.WaitGroup
	Workers      []*Worker
	Errors       []error
}

type Worker struct {
	Data interface{}
	Fnc  func(t *Worker) error
}

func (r *Processor) SetData(data interface{}) bool {
	r.errLock.Lock()
	defer r.errLock.Unlock()
	return len(r.Errors) > 0
}

func (r *Processor) AddWorker(t *Worker) {
	r.Workers = append(r.Workers, t)
}

func (r *Processor) HasErrors() bool {
	r.errLock.Lock()
	defer r.errLock.Unlock()
	return len(r.Errors) > 0
}

func (r *Processor) Run(concurrents int, tracker func(progress int)) error {
	if r.errLock == nil {
		r.errLock = &sync.Mutex{}
	}
	if r.progressLock == nil {
		r.progressLock = &sync.Mutex{}
	}
	count := len(r.Workers)
	ticks := 0
	done := make(chan int, concurrents)
	for _, worker := range r.Workers {
		r.wg.Add(1)
		go func(t *Worker, d chan int, w *sync.WaitGroup) {
			d <- 1
			err := t.Fnc(t)
			if err != nil {
				r.errLock.Lock()
				r.Errors = append(r.Errors, err)
				r.errLock.Unlock()
			}
			<-d
			w.Done()
			r.progressLock.Lock()
			ticks++
			advance := (ticks * 100) / count
			if r.progress != advance {
				r.progress = advance
				if tracker != nil {
					tracker(r.progress)
				}
			}
			r.progressLock.Unlock()
		}(worker, done, &r.wg)
	}
	r.wg.Wait()
	if len(r.Errors) > 0 {
		return ErrWorkersReturnWithError
	}
	return nil
}

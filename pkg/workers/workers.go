package workers

import (
	"sync"
	"time"
)

type WorkerPool struct {
	maxWorker  int
	queuedTask chan func()
	waitGroup  sync.WaitGroup
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		wp.waitGroup.Add(1)
		go func(workerID int) {
			for task := range wp.queuedTask {
				time.Sleep(time.Millisecond)
				task()
			}
			wp.waitGroup.Done()
		}(i + 1)
	}
}

func (wp *WorkerPool) Close() {
	close(wp.queuedTask)
}

func (wp *WorkerPool) AddTask(task func()) {
	wp.queuedTask <- task
}

func NewWorkerPool(maxWorker int) WorkerPool {
	queuedTasks := make(chan func())
	return WorkerPool{
		maxWorker:  maxWorker,
		queuedTask: queuedTasks,
		waitGroup:  sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Wait() {
	wp.waitGroup.Wait()
}

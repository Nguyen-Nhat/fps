package workers

import (
	"sync"
	"time"
)

type workerPool struct {
	maxWorker  int
	queuedTask chan func()
	waitGroup  sync.WaitGroup
}

func (wp *workerPool) Run() {
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

func (wp *workerPool) Close() {
	close(wp.queuedTask)
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTask <- task
}

func NewWorkerPool(maxWorker int) workerPool {
	queuedTasks := make(chan func())
	return workerPool{
		maxWorker:  maxWorker,
		queuedTask: queuedTasks,
		waitGroup:  sync.WaitGroup{},
	}
}

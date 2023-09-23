package boundedparallelism

import "sync"

type DigesterFunction[T any] func(args T)

type InitChannel func()

type (
	IBoundedParallelism[T any] interface {
		Execute(initChannel InitChannel, args T)
	}

	BoundedParallelism[T any] struct {
		numDigesters int
		digester     DigesterFunction[T]
	}
)

func (b *BoundedParallelism[T]) Execute(sendToChannel InitChannel, args T) {
	// 1. Create a go routine sending task to channel
	go sendToChannel()

	// 2. Create go routines for implement tasks
	var wg sync.WaitGroup
	wg.Add(b.numDigesters)
	for i := 0; i < b.numDigesters; i++ {
		go func() {
			b.digester(args)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
	}()
}

func NewBoundedParallelism[T any](numDigesters int, digester func(args T)) *BoundedParallelism[T] {
	if numDigesters == 0 {
		panic("NewBoundedParallelism but numDigesters=0")
	}
	return &BoundedParallelism[T]{numDigesters: numDigesters, digester: digester}
}

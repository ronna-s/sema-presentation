package sema4_test

import (
	"context"
	"sync"
	"testing"

	. "github.com/ronna-s/sema-presentation/request"
	"golang.org/x/sync/semaphore"
)

func process(maxHandlers int, reqs []Request) {
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(maxHandlers))

	for _, r := range reqs {
		wg.Add(1)
		sem.Acquire(context.Background(), 1)
		go func(r Request) {
			Handle(r)
			sem.Release(1)
			wg.Done()
		}(r)
	}
	wg.Wait()
}


func BenchmarkServe10(b *testing.B) {
	Serve(b.N, 10, process)
}

func BenchmarkServe100(b *testing.B) {
	Serve(b.N, 100, process)
}

func BenchmarkServe1000(b *testing.B) {
	Serve(b.N, 1000, process)
}

func BenchmarkServe10000(b *testing.B) {
	Serve(b.N, 10000, process)
}

func BenchmarkServe100000(b *testing.B) {
	Serve(b.N, 100000, process)
}

func BenchmarkServe1000000(b *testing.B) {
	Serve(b.N, 1000000, process)
}

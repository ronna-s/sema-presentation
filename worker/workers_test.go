package worker_test

import (
	"sync"
	"testing"

	. "github.com/ronna-s/sema-presentation/request"
)

func process(maxHandlers int, reqs []Request) {
	var wg sync.WaitGroup
	wg.Add(maxHandlers)

	ch := make(chan Request, 10*maxHandlers) //a good number to toy with

	for i := 0; i < maxHandlers; i++ {
		go func() {
			for r := range ch {
				r.Handle()
			}
			wg.Done()
		}()
	}
	for _, r := range reqs {
		ch <- r
	}
	close(ch)
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

package sema1_test

import (
	"sync"
	"testing"

	. "github.com/ronna-s/sema-presentation/request"
)

func process(maxHandlers int, reqs []Request) {
	var wg sync.WaitGroup
	sema := make(chan struct{}, maxHandlers)
	for _, r := range reqs {
		wg.Add(1)
		go func(r Request) {
			sema <- struct{}{}
			r.Handle()
			<-sema
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

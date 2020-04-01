package worker_test

import (
	"testing"
	"sync"
	"time"
	"sync/atomic"
)

type request struct{}

var count int64

func handle(r request) {
	atomic.AddInt64(&count, 1)
	time.Sleep(time.Microsecond * 100)
}

func process(nGoroutines int, reqs []request) {
	var wg sync.WaitGroup
	wg.Add(nGoroutines)

	ch := make(chan request, 10*nGoroutines) //a good number to toy with

	for i := 0; i < nGoroutines; i++ {
		go func() {
			for r := range ch {
				handle(r)
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

func benchmarkServe(b *testing.B, n int) {
	count = 0
	reqs := make([]request, b.N)
	process(n, reqs)
	if int(count) != b.N {
		b.Errorf("number of messages handled doesn't match, wanted: '%d' but received: '%d'", b.N, count)

	}
}

func BenchmarkServe10(b *testing.B) {
	benchmarkServe(b, 10)
}

func BenchmarkServe100(b *testing.B) {
	benchmarkServe(b, 100)
}

func BenchmarkServe1000(b *testing.B) {
	benchmarkServe(b, 1000)
}

func BenchmarkServe10000(b *testing.B) {
	benchmarkServe(b, 10000)
}

func BenchmarkServe100000(b *testing.B) {
	benchmarkServe(b, 100000)
}

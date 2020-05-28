package request

import (
	"fmt"
	"math/rand"
	"sync/atomic"
)

var nProcessed int64

type Request struct {
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func Handle(r Request) {
	atomic.AddInt64(&nProcessed, 1)
	Fib(rand.Intn(10) + 2)
}

func Serve(numMessages, n int, process func(maxHandlers int, reqs []Request)) error {
	nProcessed = 0
	reqs := make([]Request, numMessages)
	process(n, reqs)
	if int(nProcessed) != numMessages {
		return fmt.Errorf("number of messages handled doesn't match, wanted: '%d' but received: '%d'", numMessages, nProcessed)
	}
	return nil
}

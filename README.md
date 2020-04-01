# sema-presentation


#### sema1:
Take the Effective Go approach. Create a goroutine and aquire the semaphore inside the goroutine's execution.
Issue: Can run out of goroutines (depending on the rate of the incoming requests vs. the speed of `handle`.
```go
sema := make(chan struct{}, maxHandlers)
go func(r request) {
    sema <- struct{}{}
    handle(r)
    <-sema
}(r)

```
#### sema2:
Try to improve sema1 by creating the goroutine after aquiring the semaphore, the maximum
Possible issue: Slow?
```go
sema := make(chan struct{}, maxHandlers)
sema <- struct{}{}
go func(r request) {
    handle(r)
    <-sema
}(r)
```

```bash
go test github.com/ronna-s/sema-presentation/sema1 -bench=. -benchtime=1000000x
#-benchtime=Nx (N number, suffix x) means run N times, we test multiple number of workers
# against 1 million requests.
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/sema1
BenchmarkServe10-4       	 1000000	      1466 ns/op
BenchmarkServe100-4      	 1000000	       649 ns/op
BenchmarkServe1000-4     	 1000000	       641 ns/op
BenchmarkServe10000-4    	 1000000	       668 ns/op
BenchmarkServe100000-4   	 1000000	       659 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema1	4.239s
```

Quick note: ns/op in our case is actually ns per request (not operation), because in the benchmark tests we are using b.N for the number of messages.
Quick note 2: 
BenchmarkServe10 = Semaphore of 10
BenchmarkServe100 = Semaphore of 100
BenchmarkServe1000 = Semaphore of 1000
etc...

```bash
go test github.com/ronna-s/sema-presentation/sema2 -bench=. -benchtime=1000000x
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/sema2
BenchmarkServe10-4       	 1000000	       591 ns/op
BenchmarkServe100-4      	 1000000	       489 ns/op
BenchmarkServe1000-4     	 1000000	       533 ns/op
BenchmarkServe10000-4    	 1000000	       591 ns/op
BenchmarkServe100000-4   	 1000000	       609 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema2	2.832s
```

- Conclusion: Seems to be similar, an advantage to sema2 on lower numbers of goroutines (workers)

#### sema3:
Use the [Go semaphore package](https://godoc.org/golang.org/x/sync/semaphore) similarly to sema1 (aquire the semaphore inside the goroutine)
```go
go func(r request) {
    sem.Acquire(context.Background(), 1)
    handle(r)
    sem.Release(1)
}()
```
#### sema4:
Same as sema3 with the same change we made between 1 and 2 (aquire the semaphore before creating the goroutine)
```go
sem.Acquire(context.Background(), 1)
go func(r request) {
    handle(r)
    sem.Release(1)
}()

```bash
go test github.com/ronna-s/sema-presentation/sema3 -bench=. -benchtime=1000000x
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/sema3
BenchmarkServe10-4       	 1000000	      5581 ns/op
BenchmarkServe100-4      	 1000000	       526 ns/op
BenchmarkServe1000-4     	 1000000	       534 ns/op
BenchmarkServe10000-4    	 1000000	       496 ns/op
BenchmarkServe100000-4   	 1000000	       471 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema3	9.233s
```

- Conclusion, worse than sema1, similar in issue of sema2 with small amount of workers
- Note: when run with anyting higher than 3 million requests sema3 performs **_very badly_** in comparison to others. 
 
```bash
go test github.com/ronna-s/sema-presentation/sema4 -bench=. -benchtime=1000000x
```
Possible result:
```
go test github.com/ronna-s/sema-presentation/sema4 -bench=. -benchtime=1000000x
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/sema4
BenchmarkServe10-4       	 1000000	       649 ns/op
BenchmarkServe100-4      	 1000000	       433 ns/op
BenchmarkServe1000-4     	 1000000	       496 ns/op
BenchmarkServe10000-4    	 1000000	       583 ns/op
BenchmarkServe100000-4   	 1000000	       598 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema4	2.775s
```

- Conclusion: sema4 is very similar to 2 in performance (looks a little better here, but it varies).

#### worker:
Instead of a semaphore approach, why not start a finite set of go routines and share the requests between them using a queue?

```go
ch := make(chan request, 10*maxHandlers)

for i := 0; i < maxHandlers; i++ {
    go func() {
        for r := range ch {
            handle(r)
        }
		}()
	}
	for _, r := range reqs {
		ch <- r
	}
```

```bash
go test github.com/ronna-s/sema-presentation/worker -bench=. -benchtime=1000000x
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/worker
BenchmarkServe10-4       	 1000000	     13899 ns/op
BenchmarkServe100-4      	 1000000	      1253 ns/op
BenchmarkServe1000-4     	 1000000	       311 ns/op
BenchmarkServe10000-4    	 1000000	       407 ns/op
BenchmarkServe100000-4   	 1000000	       811 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/worker	17.030s
```

- Conclusion: Not amazing. Homework: profile it to see why.

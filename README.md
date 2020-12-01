# sema-presentation


#### sema1:
Take the Effective Go approach. Create a goroutine and aquire the semaphore inside the goroutine's execution.
```go
sema := make(chan struct{}, maxHandlers)
go func(r Request) {
    sema <- struct{}{}
    Handle(r)
    <-sema
}(r)

```
#### sema2:
Try to improve sema1 by creating the goroutine after aquiring the semaphore, the maximum
Possible issue: Slow?
```go
sema := make(chan struct{}, maxHandlers)
sema <- struct{}{}
go func(r Request) {
    Handle(r)
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
BenchmarkServe10-4        	 1000000	       547 ns/op
BenchmarkServe100-4       	 1000000	       474 ns/op
BenchmarkServe1000-4      	 1000000	       405 ns/op
BenchmarkServe10000-4     	 1000000	       410 ns/op
BenchmarkServe100000-4    	 1000000	       409 ns/op
BenchmarkServe1000000-4   	 1000000	       419 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema1	2.970s
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
BenchmarkServe10-4        	 1000000	       460 ns/op
BenchmarkServe100-4       	 1000000	       440 ns/op
BenchmarkServe1000-4      	 1000000	       404 ns/op
BenchmarkServe10000-4     	 1000000	       413 ns/op
BenchmarkServe100000-4    	 1000000	       415 ns/op
BenchmarkServe1000000-4   	 1000000	       415 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema2	2.710s
```

- Conclusion: Seems to be similar, an advantage to sema2 on lower numbers of goroutines (workers)

#### sema3:
Use the [Go semaphore package](https://godoc.org/golang.org/x/sync/semaphore) similarly to sema1 (aquire the semaphore inside the goroutine)
```go
go func(r Request) {
    sem.Acquire(context.Background(), 1)
    Handle(r)
    sem.Release(1)
}()
```
#### sema4:
Same as sema3 with the same change we made between 1 and 2 (aquire the semaphore before creating the goroutine)
```go
sem.Acquire(context.Background(), 1)
go func(r Request) {
    Handle(r)
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
BenchmarkServe10-4        	 1000000	      3095 ns/op
BenchmarkServe100-4       	 1000000	       345 ns/op
BenchmarkServe1000-4      	 1000000	       347 ns/op
BenchmarkServe10000-4     	 1000000	       342 ns/op
BenchmarkServe100000-4    	 1000000	       337 ns/op
BenchmarkServe1000000-4   	 1000000	       337 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema3	6.716s
```

- Conclusion, worse than sema1, similar in issue of sema2 with small amount of workers
- Note: when run with anyting higher than 3 million requests sema3 performs **_very badly_** in comparison to others. 
 
```bash
go test github.com/ronna-s/sema-presentation/sema4 -bench=. -benchtime=1000000x
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/sema4
BenchmarkServe10-4        	 1000000	       465 ns/op
BenchmarkServe100-4       	 1000000	       373 ns/op
BenchmarkServe1000-4      	 1000000	       358 ns/op
BenchmarkServe10000-4     	 1000000	       361 ns/op
BenchmarkServe100000-4    	 1000000	       367 ns/op
BenchmarkServe1000000-4   	 1000000	       365 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/sema4	2.456s
```

- Conclusion: sema4 is very similar to 2 in performance (looks a little better here, but it varies).

#### worker:
Instead of a semaphore approach, why not start a finite set of go routines and share the requests between them using a queue?

```go
ch := make(chan Request, 10*maxHandlers) //a good number to toy with

for i := 0; i < maxHandlers; i++ {
    go func() {
        for r := range ch {
            Handle(r)
        }
		}()
	}
	for _, r := range reqs {
		ch <- r
	}
    }()
```

```bash
go test github.com/ronna-s/sema-presentation/worker -bench=. -benchtime=1000000x
```
Possible result:
```
goos: darwin
goarch: amd64
pkg: github.com/ronna-s/sema-presentation/worker
BenchmarkServe10-4        	 1000000	       315 ns/op
BenchmarkServe100-4       	 1000000	       268 ns/op
BenchmarkServe1000-4      	 1000000	       272 ns/op
BenchmarkServe10000-4     	 1000000	       462 ns/op
BenchmarkServe100000-4    	 1000000	       538 ns/op
BenchmarkServe1000000-4   	 1000000	      1279 ns/op
PASS
ok  	github.com/ronna-s/sema-presentation/worker	6.176s
```

- Conclusion: Not amazing. Homework: profile it to see why.

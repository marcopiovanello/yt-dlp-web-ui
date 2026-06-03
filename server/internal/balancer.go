package internal

import (
	"container/heap"
	"log/slog"
)

type LoadBalancer struct {
	pool Pool
	done chan *Worker
}

func NewLoadBalancer(numWorker int) *LoadBalancer {
	var pool Pool

	doneChan := make(chan *Worker)

	for i := range numWorker {
		w := &Worker{
			requests: make(chan *Process, 1),
			index:    i,
		}
		go w.Work(doneChan)
		pool = append(pool, w)

		slog.Info("spawned worker", slog.Int("index", i))
	}

	return &LoadBalancer{
		pool: pool,
		done: doneChan,
	}
}

func (b *LoadBalancer) Balance(work chan *Process) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *LoadBalancer) dispatch(req *Process) {
	w := heap.Pop(&b.pool).(*Worker)
	w.requests <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *LoadBalancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
}

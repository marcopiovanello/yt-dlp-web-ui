package internal

// Pool implements heap.Interface interface as a standard priority queue
type Pool []*Worker

func (h Pool) Len() int           { return len(h) }
func (h Pool) Less(i, j int) bool { return h[i].pending < h[j].pending }

func (h Pool) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *Pool) Push(x any) { *h = append(*h, x.(*Worker)) }

func (h *Pool) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return x
}

package presents

import (
	"container/heap"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (h PresentHeap) Len() int {
	return len(h)
}

func (h PresentHeap) Less(i, j int) bool {
	if h[i].Value == h[j].Value {
		return h[i].Size < h[j].Size
	}
	return h[i].Value > h[j].Value
}

func (h PresentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PresentHeap) Push(x interface{}) {
	*h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewPresentHeap() *PresentHeap {
	h := &PresentHeap{}
	heap.Init(h)
	return h
}

func GetNCoolestPresents(presents []Present, n int) ([]Present, error) {
	h := PresentHeap{}
	if n > len(presents) {
		return nil, fmt.Errorf("n is greater (or negative) than the number of presents len(presents)=%d, n=%d", len(presents), n)
	}
	for i := range presents {
		heap.Push(&h, presents[i])
	}

	return h[:n], nil
}

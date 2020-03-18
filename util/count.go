package util

import (
	"bytes"
	"container/heap"
	"sort"
)

// URLCounter counter for a url
type URLCounter struct {
	url   []byte
	count int
}

// CounterHeap heap for finding the top n counts
type CounterHeap []URLCounter

// Less heap method less
func (h CounterHeap) Less(i, j int) bool {
	return h[i].count < h[j].count
}

// Swap heap method swap
func (h CounterHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Len heap method len
func (h CounterHeap) Len() int {
	return len(h)
}

// Push heap method push
func (h *CounterHeap) Push(x interface{}) {
	URLCounter, ok := x.(URLCounter)
	if !ok {
		panic("invalid type")
	}
	*h = append(*h, URLCounter)
}

// Pop heap method pop
func (h *CounterHeap) Pop() interface{} {
	*h = (*h)[:len(*h)-1]
	return nil
}

// CountTopN find top n counts
func CountTopN(path string, n int, bufferSize int) []URLCounter {
	partitionList := PartitionSort(path, bufferSize) // 128KB
	sorter := NewMergeSorter(partitionList)
	defer sorter.Deconstruct()
	counterHeap := make(CounterHeap, 0, n)
	row := sorter.Next()
	if row == nil {
		return counterHeap
	}
	counter := URLCounter{
		url:   row,
		count: 1,
	}
	for {
		row = sorter.Next()
		if bytes.Compare(counter.url, row) == 0 {
			counter.count++
		} else {
			if counterHeap.Len() < n {
				heap.Push(&counterHeap, counter)
			} else {
				if counter.count > counterHeap[0].count {
					counterHeap[0] = counter
					heap.Fix(&counterHeap, 0)
				}
			}
			if row == nil {
				break
			}
			counter.url = row
			counter.count = 1
		}
	}
	sort.Slice(counterHeap, func(i, j int) bool {
		return counterHeap[i].count < counterHeap[j].count
	})
	return counterHeap
}

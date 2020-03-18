package util

import (
	"bytes"
	"container/heap"
	"sort"
)

type UrlCounter struct {
	url   []byte
	count int
}

type CounterHeap []UrlCounter

func (h CounterHeap) Less(i, j int) bool {
	return h[i].count < h[j].count
}

func (h CounterHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h CounterHeap) Len() int {
	return len(h)
}

func (h *CounterHeap) Push(x interface{}) {
	urlCounter, ok := x.(UrlCounter)
	if !ok {
		panic("invalid type")
	}
	*h = append(*h, urlCounter)
}

func (h *CounterHeap) Pop() interface{} {
	*h = (*h)[:len(*h)-1]
	return nil
}

// FindTopN find top n
func CountTopN(path string, n int, bufferSize int) []UrlCounter {
	partitionList := PartitionSort(path, bufferSize) // 128KB
	sorter := NewMergeSorter(partitionList)
	defer sorter.Deconstruct()
	counterHeap := make(CounterHeap, 0, n)
	row := sorter.Next()
	if row == nil {
		return counterHeap
	}
	counter := UrlCounter{
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

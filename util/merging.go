package util

import (
	"bytes"
	"container/heap"
)

// RowPtr row data with pointer to the partition
type RowPtr struct {
	row          []byte
	partitionPtr *Partition
}

// MergeSorter merge the sorted partitions
type MergeSorter struct {
	partitionList []Partition
	heap          []RowPtr
	nextRowBuffer []byte
}

func (m *MergeSorter) Less(i, j int) bool {
	return bytes.Compare(m.heap[i].row, m.heap[j].row) < 0
}

func (m *MergeSorter) Len() int {
	return len(m.heap)
}

// Push heap method push
func (m *MergeSorter) Push(x interface{}) {
	panic("should not push")
}

// Pop heap method pop
func (m *MergeSorter) Pop() interface{} {
	m.heap = m.heap[:len(m.heap)-1]
	return nil
}

// Swap heap method swap
func (m *MergeSorter) Swap(i, j int) {
	m.heap[i], m.heap[j] = m.heap[j], m.heap[i]
}

// NewMergeSorter create a MergeSorter
func NewMergeSorter(partitionList []Partition) MergeSorter {
	merge := MergeSorter{
		partitionList: partitionList,
		nextRowBuffer: make([]byte, 128),
	}
	for i := range merge.partitionList {
		if row := merge.partitionList[i].NextRow(); row != nil {
			rowPtr := RowPtr{
				row:          append(make([]byte, 0, len(row)), row...),
				partitionPtr: &merge.partitionList[i],
			}
			merge.heap = append(merge.heap, rowPtr)
		}
	}
	heap.Init(&merge)
	return merge
}

// Next get row from the MergeSorter, the memory of the return bytes will be reused
func (m *MergeSorter) Next() []byte {
	if len(m.heap) == 0 {
		return nil
	}
	rowPtr := m.heap[0]
	m.nextRowBuffer = append(m.nextRowBuffer[0:0], rowPtr.row...)
	nextRow := rowPtr.partitionPtr.NextRow()
	if nextRow == nil {
		heap.Remove(m, 0)
	} else {
		m.heap[0].row = append(m.heap[0].row[0:0], nextRow...)
		heap.Fix(m, 0)
	}
	return m.nextRowBuffer
}

// Deconstruct free the resource of the MergeSorter
func (m *MergeSorter) Deconstruct() {
	for i := range m.partitionList {
		m.partitionList[i].Deconstruct()
	}
}

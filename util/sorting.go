package util

import (
	"bufio"
	"bytes"
	"os"
)

// GetSortedIndex get sorted index for rows in the buffer
func GetSortedIndex(buf Buffer) []int {
	sorted := make([]int, buf.Length())
	for i := range sorted {
		sorted[i] = i
	}
	// to achieve best performence on my machine, I set the workNum = 6
	ParallelSort(sorted, func(x, y int) bool {
		b1 := buf.GetRow(x)
		b2 := buf.GetRow(y)
		return bytes.Compare(b1, b2) <= 0
	}, 6)
	return sorted
}

// PartitionSort read the data into the buffer,
// sort the data, and write the sorted data into partitions
func PartitionSort(path string, partitionSize int) []Partition {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	partitionList := make([]Partition, 0, 32)

	sortBuf := NewBuffer(partitionSize)
	ioBuf := NewBuffer(partitionSize)
	ioBufReady := make(chan struct{}, 1)
	sortBufReady := make(chan struct{}, 1)
	// ioBufReady <- struct{}{}
	sortBufReady <- struct{}{}
	go func() {
		for scanner.Scan() {
			bs := scanner.Bytes()
			if ok := ioBuf.AppendBytes(bs); !ok {
				<-sortBufReady
				ioBuf, sortBuf = sortBuf, ioBuf
				ioBufReady <- struct{}{}
				ioBuf.Reset()
				ok = ioBuf.AppendBytes(bs)
				if !ok {
					panic("too big line")
				}
			}
		}
		if ioBuf.Length() != 0 {
			<-sortBufReady
			ioBuf, sortBuf = sortBuf, ioBuf
			ioBufReady <- struct{}{}
		}
		close(ioBufReady)
	}()

	for range ioBufReady {
		sortedIdx := GetSortedIndex(sortBuf)
		par := WriteToPartition(sortBuf, sortedIdx)
		partitionList = append(partitionList, par)
		sortBufReady <- struct{}{}
	}
	return partitionList
}

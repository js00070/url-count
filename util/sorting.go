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
	buf := NewBuffer(partitionSize)
	partitionList := make([]Partition, 0, 32)
	for scanner.Scan() {
		bs := scanner.Bytes()
		if ok := buf.AppendBytes(bs); !ok {
			// sort in memory
			sortedIdx := GetSortedIndex(buf)
			// write to partition
			par := WriteToPartition(buf, sortedIdx)
			partitionList = append(partitionList, par)
			buf.Reset()
			ok = buf.AppendBytes(bs)
			if !ok {
				panic("too big line")
			}
		}
	}
	if buf.Length() != 0 {
		sortedIdx := GetSortedIndex(buf)
		// log.Printf("Write to partition %v\n", len(partitionList))
		// write to partition
		par := WriteToPartition(buf, sortedIdx)
		partitionList = append(partitionList, par)
	}
	return partitionList
}

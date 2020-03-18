package util

import (
	"bufio"
	"bytes"
	"os"
	"sort"
)

// GetSortedIndex 排序
func GetSortedIndex(buf Buffer) []int {
	sorted := make([]int, buf.Length())
	for i := range sorted {
		sorted[i] = i
	}
	sort.Slice(sorted, func(i, j int) bool {
		b1 := buf.GetRow(sorted[i])
		b2 := buf.GetRow(sorted[j])
		return bytes.Compare(b1, b2) < 0
	})
	return sorted
}

// PartitionSort 分块排序
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
		// 可以换成scanner.Bytes()
		str := scanner.Text()
		if ok := buf.AppendString(str); !ok {
			// sort in memory
			// fmt.Println(buf.Length())
			sortedIdx := GetSortedIndex(buf)
			// log.Printf("Write to partition %v\n", len(partitionList))
			// write to partition
			par := WriteToPartition(buf, sortedIdx)
			partitionList = append(partitionList, par)
			buf.Reset()
			ok = buf.AppendString(str)
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

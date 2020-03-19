package util

import (
	"bufio"
	"os"
	"sort"
	"testing"
	"url-count/generator"
)

func TestMergeSort(t *testing.T) {
	filePath := generator.GenerateUrls(10, 1000, 30)
	fp, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	strList := make([]string, 0, 1000)
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		url := scanner.Text()
		strList = append(strList, url)
	}
	sort.Strings(strList)

	partitionList := PartitionSort(filePath, 512)
	sorter := NewMergeSorter(partitionList)
	defer sorter.Deconstruct()
	var row []byte
	for i := range strList {
		row = sorter.Next()
		if row == nil {
			t.Fatalf("wrong length!\n")
		}
		if strList[i] != string(row) {
			t.Errorf("%v: %v != %v\n", i, strList[i], string(row))
		}
	}
	row = sorter.Next()
	if row != nil {
		t.Fatalf("wrong length!\n")
	}
}

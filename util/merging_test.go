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
	rowList := make([][]byte, 0)
	for {
		row := sorter.Next()
		if row == nil {
			break
		}
		rowList = append(rowList, row)
	}
	if len(rowList) != len(strList) {
		t.Fatalf("length %v != %v", len(rowList), len(strList))
	}
	for i := 1; i < len(rowList); i++ {
		if strList[i] != string(rowList[i]) {
			t.Errorf("%v: %v != %v\n", i, strList[i], string(rowList[i]))
		}
	}

}

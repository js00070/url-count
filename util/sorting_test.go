package util

import (
	"sort"
	"testing"
	"url-count/generator"
)

func TestSort(t *testing.T) {
	buf := NewBuffer(1024 * 64)
	testStr := make([]string, 1024)
	for i := range testStr {
		testStr[i] = generator.GenerateRandStr(16)
	}
	for i := range testStr {
		buf.AppendString(testStr[i])
	}
	sort.Strings(testStr)
	sortedIdx := GetSortedIndex(buf)
	for i := range testStr {
		if string(buf.GetRow(sortedIdx[i])) != testStr[i] {
			t.Errorf("%v != %v\n", string(buf.GetRow(sortedIdx[i])), testStr[i])
		} else {
			t.Logf("%v == %v\n", string(buf.GetRow(sortedIdx[i])), testStr[i])
		}
	}
}

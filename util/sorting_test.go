package util

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	buf := NewBuffer(1024)
	testStr := []string{"aaa", "aaa", "bbb", "ccc", "aaa", "ddd", "ccc", "bbb", "aaa"}
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

package util

import (
	"bytes"
	"sort"
	"strings"
	"testing"
	"url-count/generator"
)

func TestParallelSort(t *testing.T) {
	testStr := make([]string, 1024)
	for i := range testStr {
		testStr[i] = generator.GenerateRandStr(16)
	}
	idx := make([]int, 1024)
	for i := range idx {
		idx[i] = i
	}
	ParallelSort(idx, func(x, y int) bool {
		b1 := testStr[x]
		b2 := testStr[y]
		return strings.Compare(b1, b2) <= 0
	}, 4)
	newStr := make([]string, 1024)
	for i, v := range idx {
		newStr[i] = testStr[v]
	}
	sort.Strings(testStr)
	for i := range newStr {
		if newStr[i] != testStr[i] {
			t.Errorf("%v != %v\n", newStr[i], testStr[i])
		}
	}
}

func BenchmarkParallelSort(b *testing.B) {
	buf := NewBuffer(1024 * 1024 * 100) // 100M
	testStr := make([]string, 1024*1024)
	for i := range testStr {
		testStr[i] = generator.GenerateRandStr(64)
	}
	for i := range testStr {
		buf.AppendString(testStr[i])
	}
	idx := make([]int, 1024*1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for i := range idx {
			idx[i] = i
		}
		b.StartTimer()
		ParallelSort(idx, func(x, y int) bool {
			b1 := buf.GetRow(x)
			b2 := buf.GetRow(y)
			return bytes.Compare(b1, b2) <= 0
		}, 4)
	}
}

package util

import (
	"testing"
)

func TestPartition(t *testing.T) {
	buf := NewBuffer(1024)
	testStr := []string{"abcdef", "11", "00000000000", "sdasdasdassdasdasdasdasd", "1", "4545422", "sdasdasdawsd", "453465342333"}
	for i := range testStr {
		buf.AppendString(testStr[i])
	}
	par := WriteToPartition(buf, []int{0, 1, 2, 3, 4, 5, 6, 7})
	defer par.Deconstruct()
	i := 0
	for {
		row := par.NextRow()
		if row == nil {
			break
		}
		if string(row) != testStr[i] {
			t.Errorf("%v != %v\n", string(row), testStr[i])
		}
		i++
	}
	if i != len(testStr) {
		t.Errorf("%v != %v", i, len(testStr))
	}
}

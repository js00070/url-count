package util

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	buf := NewBuffer(1024)
	testStr := []string{"abcdef", "11", "00000000000", "sdasdasdassdasdasdasdasd", "1", "4545422"}
	for i := range testStr {
		buf.AppendString(testStr[i])
	}
	for i := range testStr {
		b := buf.GetRow(i)
		if string(b) != testStr[i] {
			t.Errorf("%v != %v\n", string(b), testStr[i])
		}
	}
}

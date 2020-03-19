package main

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

func main() {
	fmt.Println("Hello World")
	fmt.Println(xxhash.Sum64String("Hello World!"))
	// util.CountTopN("testdata/urls_10k_100k_30percent.dat", 10)
	bs1 := make([]byte, 0, 10)
	bs2 := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bs1 = append(bs1, bs2...)
	bs2[0] = 99
	fmt.Println(bs1)
	fmt.Println(bs2)
	ch := make(chan int, 15)
	fmt.Println(len(ch))
	fmt.Println(cap(ch))
}

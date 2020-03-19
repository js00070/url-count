package util

import (
	"bufio"
	"os"
	"testing"
	"url-count/generator"
)

func TestCountTopN(t *testing.T) {
	filePath := generator.GenerateUrls(100, 10000, 20)
	topN := CountTopN(filePath, 10000, 1024*512)
	fp, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	counter := make(map[string]int, 1000)
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		url := scanner.Text()
		cnt, exist := counter[url]
		if exist {
			counter[url] = cnt + 1
		} else {
			counter[url] = 1
		}
	}
	prev := topN[0].count
	for i := range topN {
		if counter[string(topN[i].url)] != topN[i].count {
			t.Errorf("%v %v != %v\n", string(topN[i].url), counter[string(topN[i].url)], topN[i].count)
		}
		if topN[i].count > prev {
			t.Errorf("%v %v < %v\n", i, topN[i].count, prev)
		}
	}
}

func BenchmarkCountTop100WithBufferSize100M(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 8000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*100) // buffer size 100M
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

func BenchmarkCountTop100WithBufferSize200M(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 8000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*200) // buffer size 200M
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

func BenchmarkCountTop100WithBufferSize500M(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 8000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*500) // buffer size 500M
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

func BenchmarkCountTop100In10GBWithBufferSize300M(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 160000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*300) // buffer size 300M
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

func BenchmarkCountTop100In10GBWithBufferSize500M(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 160000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*500) // buffer size 500M
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

func BenchmarkCountTop100In10GBWithBufferSize1GB(b *testing.B) {
	filePath := generator.GenerateUrls(10000, 160000000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		top100 := CountTopN(filePath, 100, 1024*1024*1024) // buffer size 1GB
		b.Logf("%v, len(top100) is %v\n", i, len(top100))
	}
}

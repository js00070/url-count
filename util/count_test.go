package util

import (
	"bufio"
	"os"
	"testing"
	"url-count/generator"
)

func TestCountTopN(t *testing.T) {
	filePath := generator.GenerateUrls(200, 100000, 30)
	topN := CountTopN(filePath, 100000, 1024*128)
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
	for i := range topN {
		if counter[string(topN[i].url)] != topN[i].count {
			t.Errorf("%v %v != %v\n", string(topN[i].url), counter[string(topN[i].url)], topN[i].count)
		}
	}
}

func BenchmarkCountTop100(b *testing.B) {
	b.StopTimer()
	filePath := generator.GenerateUrls(10000, 8000000, 10)
	b.StartTimer()
	b.ReportAllocs()
	top100 := CountTopN(filePath, 100, 1024*1024*100) // buffer size 100M
	b.Logf("len(top100) is %v\n", len(top100))
}

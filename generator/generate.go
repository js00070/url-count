package generator

import (
	"fmt"
	"math/rand"
	"os"
)

const randLetters = "abcdefghijklmnopqrstuvwxyz0123456789<>?!@##$%%^&*()ABCDEFGHIJKLMNOPQRSTUVWXYZ`"

// GenerateRandStr 生成长度为n的随机串
func GenerateRandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randLetters[rand.Intn(len(randLetters))]
	}
	return string(b)
}

// GenerateUrls 随机生成含有n个url的文件,其中有百分之percent的url是从提前生成好的m个url中随机抽取的
func GenerateUrls(m, n, percent int) string {
	fileName := fmt.Sprintf("../testdata/urls_%vk_%vk_%vpercent.dat", m/1000, n/1000, percent)
	file, err := os.Open(fileName)
	if err == nil {
		file.Close()
		return fileName
	}
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	urls := make([]string, m)
	for i := 0; i < m; i++ {
		urls[i] = GenerateRandStr(64)
	}
	for i := 0; i < n; i++ {
		if rand.Intn(100) < percent {
			_, err := file.WriteString(urls[rand.Intn(m)] + "\n")
			if err != nil {
				panic(err)
			}
		} else {
			_, err := file.WriteString(GenerateRandStr(64) + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
	return fileName
}

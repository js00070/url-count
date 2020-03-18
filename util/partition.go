package util

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Partition partition
type Partition struct {
	filePath string
	fp       *os.File
	scanner  *bufio.Scanner
}

// WriteToPartition write a buffer into partiton using sorted indexes
func WriteToPartition(buf Buffer, idxList []int) Partition {
	filepath := fmt.Sprintf("../tmp/%v.part", time.Now().UnixNano())
	fp, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(fp)
	for _, idx := range idxList {
		writer.Write(buf.GetRow(idx))
		writer.WriteByte('\n')
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	fp.Seek(0, 0)
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	return Partition{
		filePath: filepath,
		fp:       fp,
		scanner:  scanner,
	}
}

// GetPath get path of the partition file
func (p *Partition) GetPath() string {
	return p.filePath
}

// NextRow get a row from the partition
func (p *Partition) NextRow() []byte {
	if p.scanner.Scan() {
		return append([]byte{}, p.scanner.Bytes()...)
	}
	return nil
}

// Deconstruct free the resource of the partition
func (p *Partition) Deconstruct() {
	p.fp.Close()
	err := os.Remove(p.filePath)
	if err != nil {
		panic(err)
	}
}

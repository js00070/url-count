package util

// Buffer memory buffer
type Buffer struct {
	len     int
	offsets []uint32
	data    []byte
}

// NewBuffer create a memory buffer with capacity
func NewBuffer(cap int) Buffer {
	if cap > 1024*1024*1024*2 {
		panic("buffer size too big")
	}
	return Buffer{
		len:     0,
		offsets: make([]uint32, 1, cap/64),
		data:    make([]byte, 0, cap),
	}
}

// GetRow get a row data by index
func (b *Buffer) GetRow(idx int) []byte {
	return b.data[b.offsets[idx]:b.offsets[(idx+1)]]
}

// AppendString append a string
func (b *Buffer) AppendString(str string) bool {
	if int(b.offsets[b.len])+len(str) > cap(b.data) {
		return false
	}
	b.data = append(b.data, str...)
	b.len++
	b.offsets = append(b.offsets, uint32(len(b.data)))
	return true
}

// AppendBytes append bytes
func (b *Buffer) AppendBytes(bs []byte) bool {
	if int(b.offsets[b.len])+len(bs) > cap(b.data) {
		return false
	}
	b.data = append(b.data, bs...)
	b.len++
	b.offsets = append(b.offsets, uint32(len(b.data)))
	return true
}

// Length get length
func (b *Buffer) Length() int {
	return b.len
}

// Reset reset memory buffer
func (b *Buffer) Reset() {
	b.len = 0
	b.offsets = b.offsets[:1]
	b.data = b.data[:0]
}

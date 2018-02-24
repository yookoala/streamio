package streamio

import (
	"io"
	"math/rand"
	"time"
)

// MockReader reads a never ending stream
// of byte with random \n in between
// until capcity (c) is met
type MockReader struct {
	b       byte
	l       int
	c       *int
	randSrc rand.Source
}

// NewMockReader returns a new MockReader that
// will write a never ending stream of same byte
// to maximum l (or the end of given byte slice)
func NewMockReader(b byte, l, c uint32) MockReader {
	ic := int(c)
	randSrc := rand.NewSource(int64(time.Now().Nanosecond()))
	return MockReader{b, int(l), &ic, randSrc}
}

// randInt returns a random positive int in [0,n).
// Expects n to be positive int.
func (nr MockReader) randIntn(in int) int {
	n := int32(in)
	if n&(n-1) == 0 { // n is power of two, can mask
		return int(int32(nr.randSrc.Int63()) & (n - 1))
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := int32(nr.randSrc.Int63())
	for v > max {
		v = int32(nr.randSrc.Int63())
	}
	return int(v % n)
}

// Read implements io.Reader
func (nr MockReader) Read(out []byte) (i int, err error) {
	if *nr.c <= 0 {
		// early exit if capicity met
		err = io.EOF
		return
	}

	// read until:
	// 1. capicty (c) met or;
	// 2. read length (l) met
	n := nr.randIntn(nr.l)
	for ; i < nr.l && i < len(out) && i < *nr.c; i++ {
		if i == n {
			out[i] = '\n'
			n += nr.randIntn(nr.l)
			continue
		}
		out[i] = nr.b
	}
	*nr.c -= i
	return
}

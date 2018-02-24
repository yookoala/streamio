package streamio_test

import (
	"bufio"
	"io"
	"testing"

	"github.com/yookoala/streamio"
)

func TestLineReader_ReadLine(t *testing.T) {

	var r io.Reader
	var lr *streamio.LineReader
	var length int

	limit := 600
	r = streamio.NewMockReader('a', 30, uint32(limit))

	lr = streamio.NewLineReader(r)
	for i := 1; i <= 3; i++ {
		line, err := lr.ReadLine()
		if err != nil {
			break
		}
		t.Logf("reader 1: %s", line)
		length += len(line) + 1
	}

	lr = streamio.NewLineReader(r)
	for {
		line, err := lr.ReadLine()
		if err != nil {
			break
		}
		t.Logf("reader 2: %s", line)
		length += len(line) + 1
	}

	// either the last line have the \n at the end, or not
	if have, want1, want2 := length, limit, limit+1; want1 != have && want2 != have {
		t.Errorf("expected %d or %d, got %d", want1, want2, have)
	}
}

func TestLineReader_Read(t *testing.T) {

	var r io.Reader
	var lr *streamio.LineReader
	var length, n int
	var err error

	limit := 600
	r = streamio.NewMockReader('a', 30, uint32(limit))

	lr = streamio.NewLineReader(r)
	for i := 1; i <= 3 && err == nil; i++ {
		line := make([]byte, 256)
		n, err = lr.Read(line)
		t.Logf("reader 1: %s", line)
		length += n
	}

	lr = streamio.NewLineReader(r)
	for err == nil {
		line := make([]byte, 256)
		n, err = lr.Read(line)
		t.Logf("reader 2: %s", line)
		length += n
	}

	// either the last line have the \n at the end, or not
	if have, want := length, limit; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
}

func BenchmarkBufioScanner(b *testing.B) {
	r := streamio.NewMockReader('a', 30, uint32(b.N))
	s := bufio.NewScanner(r)
	for s.Scan() {
		_ = s.Text() // mock read and discard
	}
}

func BenchmarkLineReader_ReadLine(b *testing.B) {
	r := streamio.NewMockReader('a', 30, uint32(b.N))
	lr := streamio.NewLineReader(r)
	var err error
	for err == nil {
		_, err = lr.ReadLine()
	}
}

func BenchmarkLineReader_Read(b *testing.B) {
	r := streamio.NewMockReader('a', 30, uint32(b.N))
	lr := streamio.NewLineReader(r)
	line := make([]byte, 4096)
	var err error
	for err == nil {
		_, err = lr.Read(line)
	}
}

package streamio

import (
	"io"
	"strings"
)

// LineReader reads lines from the supposed string reader.
//
// Different than bufio.Scanner, streamio.LineReader
// do not read ahead and cache. Hence the underlying
// io.Reader retains the read status up to the line
// that this reader has return.
//
// Which means if you stop at a specific line then
// use the underlying io.Reader elsewhere, you'd get
// contents immediately after that line.
type LineReader struct {
	r       io.Reader
	byteBuf []byte
	ended   bool
}

// Read implements io.Reader. Every read ends when
// a line end (or when the byte slice end). If
// The read ends with a '\n' character, the character
// will also be written to output slice.
func (lr LineReader) Read(out []byte) (n int, err error) {
	l := len(out)
	_, err = lr.r.Read(lr.byteBuf)
	for err == nil && n < l {
		out[n] = lr.byteBuf[0]
		n++
		if lr.byteBuf[0] == '\n' {
			break
		}
		_, err = lr.r.Read(lr.byteBuf)
	}
	return
}

// ReadLine reads a line until reach \n or io.EOF
// then return a line (striped EOL). Will return io.EOF
// if read ended
func (lr LineReader) ReadLine() (line string, err error) {
	lineBuf := make([]byte, 0, 1024)
	for _, err = lr.r.Read(lr.byteBuf); err == nil && lr.byteBuf[0] != '\n'; _, err = lr.r.Read(lr.byteBuf) {
		lineBuf = append(lineBuf, lr.byteBuf[0])
	}
	// reset length and get the line
	line = strings.TrimRight(string(lineBuf), "\n\r")
	if err == io.EOF && len(line) != 0 {
		err = nil
	}
	return
}

// NewLineReader creates a LineReader out of a supposed reader
func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{
		r:       r,
		byteBuf: make([]byte, 1),
	}
}

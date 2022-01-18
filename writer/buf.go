package writer

import (
	"bufio"
	"io"
)

const (
	defaultBufSize int = 400 * 1024
	maxBufSize     int = 1000 * 1024
)

type BufWriter struct {
	w  Writer
	bw *bufio.Writer
}

func NewBufWriter(size int, w Writer) (*BufWriter, error) {
	var bs int
	if size >= maxBufSize {
		bs = maxBufSize
	} else if size <= 0 {
		bs = defaultBufSize
	} else {
		bs = size
	}

	bw := bufio.NewWriterSize(w, bs)
	return &BufWriter{
		w:  w,
		bw: bw,
	}, nil
}

func (b *BufWriter) Write(p []byte) (n int, err error) {
	return b.bw.Write(p)
}

func (b *BufWriter) Flush() error {
	if err := b.bw.Flush(); err != nil {
		return err
	}
	if err := b.w.Flush(); err != nil {
		return err
	}

	return nil
}

func (b *BufWriter) Close() error {
	if err := b.w.Close(); err != nil {
		return err
	}

	return nil
}

func (b *BufWriter) WriteByte(c byte) error {
	return b.bw.WriteByte(c)
}

func (b *BufWriter) WriteString(s string) (int, error) {
	return b.bw.WriteString(s)
}

func (b *BufWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return b.bw.ReadFrom(r)
}

func (b *BufWriter) WriteRune(r rune) (size int, err error) {
	return b.bw.WriteRune(r)
}

func (b *BufWriter) Available() int {
	return b.bw.Available()
}

func (b *BufWriter) Buffered() int {
	return b.bw.Buffered()
}

func (b *BufWriter) Size() int {
	return b.bw.Size()
}

func (b *BufWriter) Reset() {
	b.bw.Reset(b.w)
}

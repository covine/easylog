package easylog

import (
	"bufio"
	"io"
	"sync"
)

type BufWriter interface {
	io.Writer

	Flush() error
}

type SerialBufWriter struct {
	sync.Mutex
	w *bufio.Writer
}

func NewSerialBufWriter(w io.Writer, size int) *SerialBufWriter {
	if size > 0 {
		return &SerialBufWriter{
			w: bufio.NewWriterSize(w, size),
		}
	}

	return &SerialBufWriter{
		w: bufio.NewWriter(w),
	}
}

func (s *SerialBufWriter) Write(p []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()

	return s.w.Write(p)
}

func (s *SerialBufWriter) Flush() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Flush()
}

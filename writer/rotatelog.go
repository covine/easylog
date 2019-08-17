package writer

import (
	"bufio"
	"sync"
	"time"

	rl "github.com/lestrrat-go/file-rotatelogs"
)

type RotateLogWriter struct {
	mu        sync.RWMutex
	bufWriter *bufio.Writer
	rlw       *rl.RotateLogs
}

func NewRotateLogWriter(format, linkFile string, rotate, maxAge time.Duration, bufSize int) (*RotateLogWriter, error) {
	rlw, err := rl.New(
		format,
		rl.WithLinkName(linkFile),
		rl.WithRotationTime(rotate),
		rl.WithMaxAge(maxAge),
	)
	if err != nil {
		return nil, err
	}

	var bs int
	if bufSize >= maxBufSize {
		bs = maxBufSize
	} else if bufSize <= 0 {
		bs = defaultBufSize
	} else {
		bs = bufSize
	}

	return &RotateLogWriter{
		bufWriter: bufio.NewWriterSize(rlw, bs),
		rlw:       rlw,
	}, nil
}

func (r *RotateLogWriter) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bufWriter.Write(p)
}

func (r *RotateLogWriter) Flush() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bufWriter.Flush()
}

func (r *RotateLogWriter) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = r.bufWriter.Flush()
	_ = r.rlw.Close()
}

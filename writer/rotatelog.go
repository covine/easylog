package writer

import (
	"bufio"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
)

type RotateLogWriter struct {
	fileWriter *rotatelogs.RotateLogs
	bufWriter  *bufio.Writer
}

func NewRotateLogWriter(format, linkFile string, rotate, maxAge time.Duration, bufSize int) (*RotateLogWriter, error) {
	f, err := rotatelogs.New(
		format,
		rotatelogs.WithLinkName(linkFile),
		rotatelogs.WithRotationTime(rotate),
		rotatelogs.WithMaxAge(maxAge),
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

	bufWriter := bufio.NewWriterSize(f, bs)

	return &RotateLogWriter{
		fileWriter: f,
		bufWriter:  bufWriter,
	}, nil
}

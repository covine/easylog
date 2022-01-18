package writer

import (
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
)

type RotateLogsWriter struct {
	rl *rotatelogs.RotateLogs
}

func NewRotateLogsWriter(format, linkName string, rotate, maxAge time.Duration) (*RotateLogsWriter, error) {
	rl, err := rotatelogs.New(
		format,
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithRotationTime(rotate),
		rotatelogs.WithMaxAge(maxAge),
	)
	if err != nil {
		return nil, err
	}

	return &RotateLogsWriter{
		rl: rl,
	}, nil
}

func (r *RotateLogsWriter) Write(p []byte) (n int, err error) {
	return r.rl.Write(p)
}

func (r *RotateLogsWriter) Flush() error {
	return nil
}

func (r *RotateLogsWriter) Close() error {
	return r.rl.Close()
}

func (r *RotateLogsWriter) Rotate() error {
	return r.rl.Rotate()
}

func (r *RotateLogsWriter) FileName() string {
	return r.rl.CurrentFileName()
}

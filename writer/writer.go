package writer

import "io"

type Writer interface {
	io.Writer
	io.Closer

	Flush() error
}

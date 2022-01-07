package easylog

import "io"

type SyncWriter interface {
	io.Writer
	Sync() error
}

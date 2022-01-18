package easylog

import (
	"strconv"
)

type Bytes struct {
	bytes []byte
}

func (b *Bytes) AppendByte(s byte) {
	b.bytes = append(b.bytes, s)
}

func (b *Bytes) AppendInt(i int64) {
	b.bytes = strconv.AppendInt(b.bytes, i, 10)
}

func (b *Bytes) AppendString(s string) {
	b.bytes = append(b.bytes, s...)
}

func (b *Bytes) String() string {
	return string(b.bytes)
}

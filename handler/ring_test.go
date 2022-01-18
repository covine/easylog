package handler

import (
	"testing"
	"time"

	"github.com/covine/easylog"
	"github.com/covine/easylog/diode"
	"github.com/covine/easylog/writer"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRingBufferHandler(b *testing.B) {
	defer easylog.Flush()
	defer easylog.Close()

	w, err := writer.NewBufWriter(0, writer.NewStdoutWriter())
	assert.Nil(b, err)

	rh := NewRingBufferHandler(w, StdFormatter, 165536, diode.AlertFunc(func(int) {}), 100*time.Millisecond)

	easylog.AddHandler(rh)
	defer easylog.RemoveHandler(rh)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			easylog.Info().Logf("buf stdout test")
		}
	})
}

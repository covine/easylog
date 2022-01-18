package handler

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/covine/easylog"
)

func TestStderrHandlerWithLineFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.SetLevel(easylog.DEBUG)

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStderrHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stderr")

	easylog.Info().Log()
	easylog.Info().Logf("stderr")

	t.Run("parallel info", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				easylog.Info().Logf("test parallel")
			}(i)
		}
		w.Wait()
	})

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")

	logger.EnableCaller(easylog.PANIC)
	logger.EnableStack(easylog.PANIC)
	assert.Panics(t, func() {
		logger.Panic().Logf("test panic")
	})
}

func TestStderrHandlerWithJsonFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStderrHandler(JsonFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stderr")

	easylog.Info().Log()
	easylog.Info().Logf("stderr")
	t.Run("parallel info", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				easylog.Info().Logf("test parallel")
			}(i)
		}
		w.Wait()
	})

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")

	logger.EnableCaller(easylog.PANIC)
	logger.EnableStack(easylog.PANIC)
	assert.Panics(t, func() {
		logger.Panic().Logf("test panic")
	})
}

func TestStdoutHandlerWithLineFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStdoutHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}

func TestStdoutHandlerWithJsonFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStdoutHandler(JsonFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}

func BenchmarkStderrHandler(b *testing.B) {
	defer easylog.Flush()

	h := NewStderrHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			easylog.Info().Logf("stderr test")
		}
	})
}

func BenchmarkBufStderrHandler(b *testing.B) {
	defer easylog.Flush()

	h, err := NewBufStderrHandler(StdFormatter)
	assert.Nil(b, err)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			easylog.Info().Logf("buf stderr test")
		}
	})
}

func BenchmarkStdoutHandler(b *testing.B) {
	defer easylog.Flush()

	h := NewStdoutHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			easylog.Info().Logf("stdout test")
		}
	})
}

func BenchmarkBufStdoutHandler(b *testing.B) {
	defer easylog.Flush()

	h, err := NewBufStdoutHandler(StdFormatter)
	assert.Nil(b, err)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			easylog.Info().Logf("buf stdout test")
		}
	})
}

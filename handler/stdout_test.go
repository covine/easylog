package handler

import (
	"sync"
	"testing"

	"github.com/govine/easylog"
)

func TestStdout(t *testing.T) {

	t.Run("logger", func(t *testing.T) {
		l := easylog.GetLogger("test")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		stdHandler := NewStdoutHandler(format)
		l.AddHandler(stdHandler)

		l.Debug().Msg("debug: %s", "test")
		l.Flush()
		l.Info().Msg("info: %s", "test")
		l.Warn().Msg("warn: %s", "test")
		l.Warning().Msg("warning: %s", "test")
		l.Error().Msg("error: %s", "test")
		l.Fatal().Msg("fatal: %s", "test")
		l.Flush()

		var w sync.WaitGroup
		for i := 0; i < 10000; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				s := easylog.NewCachedLogger(l)
				s.SetLevel(easylog.DEBUG)
				s.SetPropagate(true)
				s.EnableFrame(easylog.DEBUG)
				s.EnableFrame(easylog.INFO)
				s.EnableFrame(easylog.ERROR)
				s.EnableFrame(easylog.WARN)
				s.EnableFrame(easylog.FATAL)
				defer func() {
					s.Flush()
					s.Close()
				}()

				s.Debug().Msg("s debug: %s", "test")
				s.Info().Msg("s info: %s", "test")
				s.Warn().Msg("s warn: %s", "test")
				s.Warning().Msg("s warning: %s", "test")
				s.Error().Msg("s error: %s", "test")
				s.Fatal().Msg("s fatal: %s", "test")
			}(i)
		}
		w.Wait()

		l.Flush()
		l.Close()
	})
}

package handler

import (
	"sync"
	"testing"

	"git.qutoutiao.net/govine/easylog"
)

func TestStdout(t *testing.T) {

	t.Run("logger", func(t *testing.T) {
		l := easylog.GetLogger("test")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		stdHandler := NewStdoutHandler(format)
		l.AddHandler(stdHandler)

		l.Debug("debug: %s", "test")
		l.Flush()
		l.Info("info: %s", "test")
		l.Warn("warn: %s", "test")
		l.Warning("warning: %s", "test")
		l.Error("error: %s", "test")
		l.Fatal("fatal: %s", "test")
		l.Flush()

		var w sync.WaitGroup
		for i := 0; i < 10000; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				s := easylog.NewCachedLogger(l)
				s.SetLevel(easylog.DEBUG)
				s.SetPropagate(true)
				defer s.Flush()

				s.Debug("s debug: %s", "test")
				s.Info("s info: %s", "test")
				s.Warn("s warn: %s", "test")
				s.Warning("s warning: %s", "test")
				s.Error("s error: %s", "test")
				s.Fatal("s fatal: %s", "test")
			}(i)
		}
		w.Wait()

		l.Flush()
	})
}

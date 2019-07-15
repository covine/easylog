package handler

import (
	"sync"
	"testing"

	"git.qutoutiao.net/govine/easylog/filter"

	"git.qutoutiao.net/govine/easylog"
	"git.qutoutiao.net/govine/easylog/writer"
)

func TestLog(t *testing.T) {

	t.Run("logger", func(t *testing.T) {
		l := easylog.GetLogger("test")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		dw, err := writer.NewRotateFileWriter(30, "./log/cpc.debug", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		DebugFileHandler := easylog.NewHandler(&FileHandler{FileWriter: dw})
		DebugFileHandler.SetLevel(easylog.DEBUG)
		DebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

		fw, err := writer.NewRotateFileWriter(30, "./log/cpc.fatal", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		FatalFileHandler := easylog.NewHandler(&FileHandler{FileWriter: fw})
		FatalFileHandler.SetLevel(easylog.FATAL)
		FatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

		ww, err := writer.NewRotateFileWriter(30, "./log/cpc.warn", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		WarnFileHandler := easylog.NewHandler(&FileHandler{FileWriter: ww})
		WarnFileHandler.SetLevel(easylog.WARN)
		WarnFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.WARN})

		l.AddHandler(DebugFileHandler)
		l.AddHandler(FatalFileHandler)
		l.AddHandler(WarnFileHandler)

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

				s := easylog.GetSparkLogger()
				s.SetLevel(easylog.DEBUG)
				defer func() {
					s.Flush()
					s.Close()
				}()

				SDebugFileHandler := easylog.NewHandler(&StoreHandler{
					fileWriter: dw,
					logs:       make([]string, 0),
					flushed:    false,
				})
				SDebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

				SFatalFileHandler := easylog.NewHandler(&StoreHandler{
					fileWriter: fw,
					logs:       make([]string, 0),
					flushed:    false,
				})
				SFatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

				SWarnFileHandler := easylog.NewHandler(&StoreHandler{
					fileWriter: ww,
					logs:       make([]string, 0),
					flushed:    false,
				})
				SWarnFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.WARN})

				s.AddHandler(SDebugFileHandler)
				s.AddHandler(SFatalFileHandler)
				s.AddHandler(SWarnFileHandler)

				s.Debug("s debug: %s", "test")
				s.Info("s info: %s", "test")
				s.Warn("s warn: %s", "test")
				s.Warning("s warning: %s", "test")
				s.Error("s error: %s", "test")
				s.Fatal("s fatal: %s", "test")
			}(i)
		}
		w.Wait()

		dw.Close()
		fw.Close()
		ww.Close()

		l.Flush()
		l.Close()
	})
}

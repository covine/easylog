package handler

import (
	"testing"

	"git.qutoutiao.net/govine/easylog/store"

	"git.qutoutiao.net/govine/easylog/filter"

	"git.qutoutiao.net/govine/easylog/formatter"

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
		DebugFileHandler, err := NewFileHandler(easylog.DEBUG, dw)
		if err != nil {
			t.Error(err)
		}
		DebugFileHandler.SetFormatter(&formatter.SimpleFormatter{})
		DebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

		fw, err := writer.NewRotateFileWriter(30, "./log/cpc.fatal", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		FatalFileHandler, err := NewFileHandler(easylog.FATAL, fw)
		if err != nil {
			t.Error(err)
		}
		FatalFileHandler.SetFormatter(&formatter.SimpleFormatter{})
		FatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

		ww, err := writer.NewRotateFileWriter(30, "./log/cpc.warn", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		WarnFileHandler, err := NewFileHandler(easylog.WARN, ww)
		if err != nil {
			t.Error(err)
		}
		WarnFileHandler.SetFormatter(&formatter.SimpleFormatter{})
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

		s, err := store.NewStoreLogger()
		if err != nil {
			t.Error(err)
			return
		}
		s.SetLevel(easylog.DEBUG)

		SDebugFileHandler, err := NewStoreHandler(easylog.DEBUG, dw)
		if err != nil {
			t.Error(err)
		}
		SDebugFileHandler.SetFormatter(&formatter.SimpleFormatter{})
		SDebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

		SFatalFileHandler, err := NewStoreHandler(easylog.FATAL, fw)
		if err != nil {
			t.Error(err)
		}
		SFatalFileHandler.SetFormatter(&formatter.SimpleFormatter{})
		SFatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

		SWarnFileHandler, err := NewStoreHandler(easylog.WARN, ww)
		if err != nil {
			t.Error(err)
		}
		SWarnFileHandler.SetFormatter(&formatter.SimpleFormatter{})
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
		s.Flush()

		dw.Close()
		fw.Close()
		ww.Close()

		s.Close()
		l.Close()
	})
}

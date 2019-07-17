package handler

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"git.qutoutiao.net/govine/easylog"
	"git.qutoutiao.net/govine/easylog/filter"
	"git.qutoutiao.net/govine/easylog/writer"
)

func format(record easylog.Record) string {
	var prefix string
	switch record.Level {
	case easylog.FATAL:
		prefix = "FATAL: "
	case easylog.ERROR:
		prefix = "ERROR: "
	case easylog.WARNING:
		prefix = "WARNING: "
	case easylog.INFO:
		prefix = "NOTICE: "
	case easylog.DEBUG:
		prefix = "DEBUG: "
	default:
		prefix = "UNKNOWN LEVEL: "
	}

	var body string
	var head string
	if record.Level == easylog.INFO {
		head = prefix + " " + time.Now().Format("2006-01-02 15:04:05") + " * "
	} else {
		_, file, line, ok := runtime.Caller(8)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = path.Base(file)
		}
		head = prefix + " " + time.Now().Format("2006-01-02 15:04:05") + " " + file + " [" + strconv.Itoa(line) + "] * "
	}
	if record.Args != nil && len(record.Args) > 0 {
		body = fmt.Sprintf(record.Msg, record.Args...)
	} else {
		body = record.Msg
	}
	return head + body
}

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
		DebugFileHandler := NewWriteHandler(format, dw)
		DebugFileHandler.SetLevel(easylog.DEBUG)
		DebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

		fw, err := writer.NewRotateFileWriter(30, "./log/cpc.fatal", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		FatalFileHandler := NewWriteHandler(format, fw)
		FatalFileHandler.SetLevel(easylog.FATAL)
		FatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

		ww, err := writer.NewRotateFileWriter(30, "./log/cpc.warn", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		WarnFileHandler := NewWriteHandler(format, ww)
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

				SDebugFileHandler := NewStoreWriteHandler(format, dw)
				SDebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

				SFatalFileHandler := NewStoreWriteHandler(format, fw)
				SFatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

				SWarnFileHandler := NewStoreWriteHandler(format, ww)
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

				go s.Flush()
				go s.Close()
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

package handler

import (
	"fmt"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"git.qutoutiao.net/govine/easylog"
	"git.qutoutiao.net/govine/easylog/filter"
	"git.qutoutiao.net/govine/easylog/writer"
)

func format(record *easylog.Record) string {
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
		file := "???"
		line := 0
		if record.OK {
			file = path.Base(record.File)
			line = record.Line
		}
		head = prefix + " " + time.Now().Format("2006-01-02 15:04:05") + " " + file + " [" + strconv.Itoa(line) + "] * "
	}
	if record.Args != nil && len(record.Args) > 0 {
		body = fmt.Sprintf(record.Message, record.Args...)
	} else {
		body = record.Message
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
		l.EnableFrame(easylog.WARN)

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
				s.EnableFrame(easylog.WARN)
				s.EnableFrame(easylog.FATAL)

				s.Debug().Msg("s debug: %s", "test")
				s.Info().Msg("s info: %s", "test")
				s.Warn().Msg("s warn: %s", "test")
				s.Warning().Msg("s warning: %s", "test")
				s.Error().Msg("s error: %s", "test")
				s.Fatal().Msg("s fatal: %s", "test")

				go func() {
					s.Flush()
					s.Close()
				}()
			}(i)
		}
		w.Wait()

		l.Flush()
		l.Close()
	})
}

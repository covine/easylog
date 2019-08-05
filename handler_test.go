package easylog_test

import (
	"fmt"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/govine/easylog/filter"
	"github.com/govine/easylog/writer"

	"github.com/govine/easylog"
	"github.com/govine/easylog/handler"
)

func TestFileHandler(t *testing.T) {
	fileHandler, err := handler.NewFileHandler("./log/file.log", nil)
	if err != nil {
		t.Error(err)
	}
	easylog.AddHandler(fileHandler)
	easylog.SetLevel(easylog.DEBUG)
	easylog.Debug().Fields(map[string]interface{}{"name": "dog"}).Msg("hello")
	// Output:
	easylog.RemoveHandler(fileHandler)
}

func TestStderrHandler(t *testing.T) {
	stderrHandler := handler.NewStderrHandler(nil)
	easylog.AddHandler(stderrHandler)
	easylog.SetLevel(easylog.DEBUG)
	easylog.Debug().Msg("hello world")
	easylog.Fatal().Msg("error")
	easylog.RemoveHandler(stderrHandler)
}

func TestStderrConcurrence(t *testing.T) {

	t.Run("stderr", func(t *testing.T) {
		l := easylog.GetLogger("stderr")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		stdHandler := handler.NewStderrHandler(nil)
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
		l.RemoveHandler(stdHandler)
		l.Close()
	})
}

func TestStdoutConcurrence(t *testing.T) {

	t.Run("stdout", func(t *testing.T) {
		l := easylog.GetLogger("stdout")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		stdHandler := handler.NewStdoutHandler(nil)
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
		l.RemoveHandler(stdHandler)

		l.Close()
	})
}

func TestHasHandlers(t *testing.T) {
	root := easylog.GetRootLogger()
	if root.HasHandler() {
		t.Errorf("root has no handler when init")
		return
	}

	stdHandler := handler.NewStderrHandler(nil)
	root.AddHandler(stdHandler)

	if !root.HasHandler() {
		t.Errorf("root must has handler after add handler")
		return
	}

	root.RemoveHandler(stdHandler)
	if root.HasHandler() {
		t.Errorf("root must has no handler after remove handler")
		return
	}
}

func fFormat(record *easylog.Record) string {
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

func TestRotateFileHandler(t *testing.T) {

	t.Run("rotateFile", func(t *testing.T) {
		l := easylog.GetLogger("rotateFile")
		l.SetLevel(easylog.DEBUG)
		l.SetPropagate(false)

		dw, err := writer.NewRotateFileWriter(30, "./log/debug.log", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		DebugFileHandler := handler.NewRotateFileHandler(fFormat, dw)
		DebugFileHandler.SetLevel(easylog.DEBUG)
		DebugFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})

		fw, err := writer.NewRotateFileWriter(30, "./log/fatal.log", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		FatalFileHandler := handler.NewRotateFileHandler(fFormat, fw)
		FatalFileHandler.SetLevel(easylog.FATAL)
		FatalFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

		ww, err := writer.NewRotateFileWriter(30, "./log/warn.log", 409600)
		if err != nil {
			t.Error(err)
			return
		}
		WarnFileHandler := handler.NewRotateFileHandler(fFormat, ww)
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
					s.Close()
				}()
			}(i)
		}
		w.Wait()

		l.Flush()
		l.RemoveHandler(DebugFileHandler)
		l.RemoveHandler(FatalFileHandler)
		l.RemoveHandler(WarnFileHandler)
		l.Close()
	})
}

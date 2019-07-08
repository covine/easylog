package formatter

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"

	"git.qutoutiao.net/govine/easylog"
)

type SimpleFormatter struct {
}

func (s *SimpleFormatter) Format(record easylog.Record) string {
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
		_, file, line, ok := runtime.Caller(6)
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

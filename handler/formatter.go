package handler

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/covine/easylog"
)

type Formatter func(e *easylog.Event) ([]byte, error)

func JsonFormatter(e *easylog.Event) ([]byte, error) {
	m := make(map[string]interface{})
	m["logger"] = e.GetLogger().Name()
	if e.GetTags() != nil {
		m["tag"] = e.GetTags()
	}
	if e.GetKvs() != nil {
		m["kvs"] = e.GetKvs()
	}
	m["time"] = e.GetTime().Format("2006-01-02 15:04:05")
	m["level"] = e.GetLevel().String()
	m["caller"] = map[string]interface{}{
		"ok":   e.GetCaller().GetOK(),
		"pc":   e.GetCaller().GetPC(),
		"file": e.GetCaller().GetFile(),
		"func": e.GetCaller().GetFunc(),
		"line": e.GetCaller().GetLine(),
	}
	m["msg"] = e.GetMsg()
	m["stack"] = e.GetStack()
	m["extra"] = e.GetExtra()

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func StdFormatter(e *easylog.Event) ([]byte, error) {
	b := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(b)

	level := e.GetLevel()
	switch level {
	case easylog.DEBUG:
		buf.WriteString(level.String())
		buf.WriteString("  ")
	case easylog.INFO:
		buf.WriteString(level.String())
		buf.WriteString("   ")
	case easylog.WARN:
		buf.WriteString(Yellow)
		buf.WriteString(level.String())
		buf.WriteString("   ")
		buf.WriteString(Reset)
	case easylog.ERROR:
		fallthrough
	case easylog.PANIC:
		fallthrough
	case easylog.FATAL:
		buf.WriteString(Red)
		buf.WriteString(level.String())
		buf.WriteString("  ")
		buf.WriteString(Reset)
	default:
		buf.WriteString("UNKNOWN")
	}

	buf.WriteString(": ")
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buf.WriteString(" *")
	buf.WriteString(" logger: ")
	buf.WriteString(e.GetLogger().Name())
	buf.WriteString(" *")

	if e.GetCaller().GetOK() {
		buf.WriteString(" ")
		buf.WriteString(path.Base(e.GetCaller().GetFile()))
		buf.WriteString(" ")
		f := strings.Split(e.GetCaller().GetFunc(), ".")
		if len(f) > 0 {
			buf.WriteString(f[len(f)-1])
			buf.WriteString(" ")
		}
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(e.GetCaller().GetLine()))
		buf.WriteString("]")
		buf.WriteString(" *")
	}

	buf.WriteString(" ")
	buf.WriteString(e.GetMsg())

	if len(e.GetStack()) > 0 {
		buf.WriteString(" stack: \n")
		buf.WriteString(e.GetStack())
	}

	return buf.Bytes(), nil
}

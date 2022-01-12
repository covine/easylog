package easylog

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Event struct {
	Logger *logger

	Time  time.Time
	Level Level
	Tags  map[string]interface{}
	Kvs   map[string]interface{}
	Msg   string
	Extra interface{}

	OK   bool
	PC   uintptr
	File string
	Line int
	Func string

	Stack string
}

var _eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

func newEvent(logger *logger, level Level) *Event {
	r := _eventPool.Get().(*Event)

	r.Logger = logger

	r.Time = time.Time{}
	r.Level = level
	r.Tags = nil
	r.Kvs = nil
	r.Msg = ""
	r.Extra = nil

	r.OK = false
	r.PC = 0
	r.File = ""
	r.Line = 0
	r.Func = ""

	r.Stack = ""

	return r
}

func putEvent(r *Event) {
	_eventPool.Put(r)
}

func (e *Event) Tag(k string, v interface{}) *Event {
	if e == nil {
		return e
	}

	if e.Tags == nil {
		e.Tags = make(map[string]interface{})
	}

	e.Tags[k] = v

	return e
}

func (e *Event) Kv(k string, v interface{}) *Event {
	if e == nil {
		return e
	}

	if e.Kvs == nil {
		e.Kvs = make(map[string]interface{})
	}

	e.Kvs[k] = v

	return e
}

func (e *Event) Attach(extra interface{}) *Event {
	if e == nil {
		return e
	}

	e.Extra = extra

	return e
}

func (e *Event) Log() {
	e.log("", 2)
}

func (e *Event) Logf(msg string, args ...interface{}) {
	e.log(fmt.Sprintf(msg, args...), 2)
}

func (e *Event) log(msg string, skip int) {
	if e == nil {
		return
	}

	e.Time = time.Now()
	e.Msg = msg

	if e.Logger.logCaller(e.Level) {
		frame, ok := e.getCallerFrame(skip)
		if !ok {
			// TODO design
			// _, _ = fmt.Fprintf(e.Logger.errHandler, "[%v] [%v] [%v]:get caller failed\n",
			//e.Logger.name, e.Level, e.Time,
			//)
			_ = e.Logger.errorHandler.Flush()
		}

		e.OK = ok
		e.PC = frame.PC
		e.File = frame.File
		e.Line = frame.Line
		e.Func = frame.Function
	}
	// TODO need stack?

	e.Logger.handle(e)

	e.Logger.couldEnd(e.Level, e.Msg)
}

func (e *Event) getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	pc := make([]uintptr, 1)

	// ignore Caller(1) and getCallerFrame(2)
	n := runtime.Callers(2+skip, pc)
	if n < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()

	return frame, frame.PC != 0
}

package easylog

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var _bytesPool = &sync.Pool{
	New: func() interface{} {
		return &Bytes{
			bytes: make([]byte, 0, 1024),
		}
	},
}

func newBytes() *Bytes {
	b := _bytesPool.Get().(*Bytes)
	b.bytes = b.bytes[:0]
	return b
}

func putBytes(b *Bytes) {
	const max = 1 << 16
	if cap(b.bytes) > max {
		return
	}
	_bytesPool.Put(b)
}

type pcs struct {
	pcs []uintptr
}

var _pcsPool = &sync.Pool{
	New: func() interface{} {
		return &pcs{make([]uintptr, 64)}
	},
}

func newPcs() *pcs {
	return _pcsPool.Get().(*pcs)
}

func putPcs(p *pcs) {
	_pcsPool.Put(p)
}

type caller struct {
	ok   bool
	pc   uintptr
	file string
	line int
	fc   string
}

func (c *caller) GetOK() bool {
	return c.ok
}

func (c *caller) GetPC() uintptr {
	return c.pc
}

func (c *caller) GetFile() string {
	return c.file
}

func (c *caller) GetLine() int {
	return c.line
}

func (c *caller) GetFunc() string {
	return c.fc
}

type Event struct {
	logger *logger

	time  time.Time
	level Level
	tags  map[interface{}]interface{}
	kvs   map[interface{}]interface{}
	msg   string
	extra interface{}

	caller caller
	stack  string
}

var _eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

func newEvent(logger *logger, level Level) *Event {
	r := _eventPool.Get().(*Event)

	r.logger = logger

	r.time = time.Time{}
	r.level = level
	r.tags = nil
	r.kvs = nil
	r.msg = ""
	r.extra = nil

	r.caller.ok = false
	r.caller.pc = 0
	r.caller.file = ""
	r.caller.line = 0
	r.caller.fc = ""

	r.stack = ""

	return r
}

func (e *Event) Tag(k interface{}, v interface{}) *Event {
	if e == nil {
		return e
	}

	if e.tags == nil {
		e.tags = make(map[interface{}]interface{})
	}

	e.tags[k] = v

	return e
}

func (e *Event) Kv(k, v interface{}) *Event {
	if e == nil {
		return e
	}

	if e.kvs == nil {
		e.kvs = make(map[interface{}]interface{})
	}

	e.kvs[k] = v

	return e
}

func (e *Event) GetKvs() map[interface{}]interface{} {
	return e.kvs
}

func (e *Event) Attach(extra interface{}) *Event {
	if e == nil {
		return e
	}

	e.extra = extra

	return e
}

func (e *Event) Log() {
	e.log("", 2)
}

func (e *Event) Logf(msg string, args ...interface{}) {
	if len(args) > 0 {
		e.log(fmt.Sprintf(msg, args...), 2)
	} else {
		e.log(msg, 2)
	}
}

func (e *Event) GetLogger() *logger {
	return e.logger
}

func (e *Event) GetTime() time.Time {
	return e.time
}

func (e *Event) GetTags() map[interface{}]interface{} {
	return e.tags
}

func (e *Event) GetCaller() *caller {
	return &e.caller
}

func (e *Event) GetLevel() Level {
	return e.level
}

func (e *Event) GetStack() string {
	return e.stack
}

func (e *Event) GetMsg() string {
	return e.msg
}

func (e *Event) GetExtra() interface{} {
	return e.extra
}

func (e *Event) Clone() *Event {
	r := _eventPool.Get().(*Event)

	r.logger = e.logger

	r.time = e.time
	r.level = e.level
	r.tags = e.tags
	r.kvs = e.kvs
	r.msg = e.msg
	r.extra = e.extra

	r.caller = e.caller

	r.stack = e.stack

	return r
}

func (e *Event) Put() {
	_eventPool.Put(e)
}

func (e *Event) log(msg string, skip int) {
	if e == nil {
		return
	}

	e.time = time.Now()
	e.msg = msg

	if e.logger.logCaller(e.level) {
		frame, ok := e.getCallerFrame(skip)
		if !ok {
			_ = e.logger.errorHandler.Handle(
				errors.New(
					fmt.Sprintf("[%v] [%v] [%v]:get caller failed\n", e.logger.name, e.level, e.time),
				),
			)
		}

		e.caller.ok = ok
		e.caller.pc = frame.PC
		e.caller.file = frame.File
		e.caller.line = frame.Line
		e.caller.fc = frame.Function
	}

	if e.logger.logStack(e.level) {
		e.stack = e.stacktrace(skip)
	}

	e.logger.handle(e)

	e.logger.couldEnd(e.level, e.msg)
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

func (e *Event) stacktrace(skip int) string {
	bs := newBytes()
	defer putBytes(bs)

	p := newPcs()
	defer putPcs(p)

	var numFrames int
	for {
		numFrames = runtime.Callers(skip+2, p.pcs)
		if numFrames < len(p.pcs) {
			break
		}

		p = &pcs{pcs: make([]uintptr, len(p.pcs)*2)}
	}

	i := 0
	frames := runtime.CallersFrames(p.pcs[:numFrames])

	for {
		frame, more := frames.Next()
		if i != 0 {
			bs.AppendByte('\n')
		}
		i++
		bs.AppendByte('\t')
		bs.AppendString(frame.Function)
		bs.AppendByte('\n')
		bs.AppendByte('\t')
		bs.AppendString(frame.File)
		bs.AppendByte(':')
		bs.AppendInt(int64(frame.Line))
		if !more {
			break
		}
	}

	return bs.String()
}

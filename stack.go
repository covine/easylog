package easylog

import (
	"runtime"
	"sync"

	"github.com/covine/easylog/buffer"
)

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get

	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}
)

func takeStacktrace(skip int) string {
	buffer := Get()
	defer buffer.Free()

	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int
	for {
		// Skip the call to runtime.Callers and takeStacktrace so that the
		// program counters start at the caller of takeStacktrace.
		numFrames = runtime.Callers(skip+2, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	i := 0
	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buffer.AppendByte('\n')
		}
		i++
		buffer.AppendString(frame.Function)
		buffer.AppendByte('\n')
		buffer.AppendByte('\t')
		buffer.AppendString(frame.File)
		buffer.AppendByte(':')
		buffer.AppendInt(int64(frame.Line))
	}

	return buffer.String()
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}

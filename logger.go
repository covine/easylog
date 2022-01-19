package easylog

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	PANIC
	FATAL

	_MIN = DEBUG
	_MAX = FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger is not thread safe
// Make sure to configure the Logger before emitting logs,
// And do not reconfigure the Logger during runtime.
type Logger struct {
	manager     *manager
	parent      *Logger
	placeholder bool
	children    map[*Logger]struct{}

	name      string
	propagate bool
	level     Level

	handlers     []Handler
	errorHandler ErrorHandler

	caller map[Level]bool
	stack  map[Level]bool

	tags map[interface{}]interface{}
	kvs  map[interface{}]interface{}
}

func newLogger() *Logger {
	return &Logger{
		children:     make(map[*Logger]struct{}),
		handlers:     make([]Handler, 0),
		errorHandler: NewNopErrorHandler(),
		caller:       make(map[Level]bool),
		stack:        make(map[Level]bool),
		tags:         make(map[interface{}]interface{}),
		kvs:          make(map[interface{}]interface{}),
	}
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) SetPropagate(propagate bool) {
	l.propagate = propagate
}

func (l *Logger) GetPropagate() bool {
	return l.propagate
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) AddHandler(h Handler) {
	if h == nil {
		return
	}

	for _, handler := range l.handlers {
		if handler == h {
			return
		}
	}

	l.handlers = append(l.handlers, h)
}

func (l *Logger) RemoveHandler(h Handler) {
	if h == nil {
		return
	}

	for i, handler := range l.handlers {
		if handler == h {
			l.handlers = append(l.handlers[:i], l.handlers[i+1:]...)
			return
		}
	}
}

func (l *Logger) ResetHandler() {
	l.handlers = make([]Handler, 0)
}

func (l *Logger) SetErrorHandler(w ErrorHandler) {
	l.errorHandler = w
}

func (l *Logger) EnableCaller(level Level) {
	if level >= _MIN && level <= _MAX {
		l.caller[level] = true
	}
}

func (l *Logger) DisableCaller(level Level) {
	if level >= _MIN && level <= _MAX {
		l.caller[level] = false
	}
}

func (l *Logger) EnableStack(level Level) {
	if level >= _MIN && level <= _MAX {
		l.stack[level] = true
	}
}

func (l *Logger) DisableStack(level Level) {
	if level >= _MIN && level <= _MAX {
		l.stack[level] = false
	}
}

func (l *Logger) SetTag(k interface{}, v interface{}) {
	l.tags[k] = v
}

func (l *Logger) DelTag(k interface{}) {
	delete(l.tags, k)
}

func (l *Logger) ResetTag() {
	l.tags = make(map[interface{}]interface{})
}

func (l *Logger) Tags() map[interface{}]interface{} {
	return l.tags
}

func (l *Logger) SetKv(k interface{}, v interface{}) {
	l.kvs[k] = v
}

func (l *Logger) DelKv(k interface{}) {
	delete(l.kvs, k)
}

func (l *Logger) ResetKv() {
	l.kvs = make(map[interface{}]interface{})
}

func (l *Logger) Kvs() map[interface{}]interface{} {
	return l.kvs
}

func (l *Logger) Debug() *Event {
	return l.log(DEBUG, nil)
}

func (l *Logger) Info() *Event {
	return l.log(INFO, nil)
}

func (l *Logger) Warn() *Event {
	return l.log(WARN, nil)
}

func (l *Logger) Error() *Event {
	return l.log(ERROR, nil)
}

func (l *Logger) Panic() *Event {
	return l.log(PANIC, func(v interface{}) {
		l.Flush()
		panic(v)
	})
}

func (l *Logger) Fatal() *Event {
	return l.log(FATAL, func(v interface{}) {
		l.Flush()
		l.Close()
		exit(1)
	})
}

func (l *Logger) Flush() {
	for _, handler := range l.handlers {
		if err := handler.Flush(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}

	// ignore error produced by errorHandler
	_ = l.errorHandler.Flush()
}

func (l *Logger) Close() {
	for _, handler := range l.handlers {
		if err := handler.Close(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}

	// ignore error produced by errorHandler
	_ = l.errorHandler.Close()
}

func (l *Logger) logCaller(level Level) bool {
	if need, ok := l.caller[level]; ok {
		return need
	}

	return false
}

func (l *Logger) logStack(level Level) bool {
	if need, ok := l.stack[level]; ok {
		return need
	}

	return false
}

// couldEnd could end the Logger with panic or os.exit().
func (l *Logger) couldEnd(level Level, v interface{}) {
	// Note: If there is any level bigger than PANIC added, the logic here should be updated.
	switch level {
	case PANIC:
		l.Flush()
		panic(v)
	case FATAL:
		l.Flush()
		l.Close()
		exit(1)
	}
}

func (l *Logger) log(level Level, done func(interface{})) *Event {
	if level < l.level {
		if done != nil {
			done("")
		}
		return nil
	}

	return newEvent(l, level)
}

func (l *Logger) handle(event *Event) {
	defer event.Put()

	if event.level < l.level {
		return
	}

	for _, handler := range l.handlers {
		next, err := handler.Handle(event)
		if err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
		if !next {
			return
		}
	}

	if l.propagate && l.parent != nil {
		l.parent.handle(event)
	}
}

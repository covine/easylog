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

// logger is not thread safe
// Make sure to configure the logger before emitting logs,
// And do not reconfigure the logger during runtime.
type logger struct {
	manager     *manager
	parent      *logger
	placeholder bool
	children    map[*logger]struct{}

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

func newLogger() *logger {
	return &logger{
		children:     make(map[*logger]struct{}),
		handlers:     make([]Handler, 0),
		errorHandler: NewNopErrorHandler(),
		caller:       make(map[Level]bool),
		stack:        make(map[Level]bool),
		tags:         make(map[interface{}]interface{}),
		kvs:          make(map[interface{}]interface{}),
	}
}

func (l *logger) Name() string {
	return l.name
}

func (l *logger) SetPropagate(propagate bool) {
	l.propagate = propagate
}

func (l *logger) GetPropagate() bool {
	return l.propagate
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) GetLevel() Level {
	return l.level
}

func (l *logger) AddHandler(h Handler) {
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

func (l *logger) RemoveHandler(h Handler) {
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

func (l *logger) ResetHandler() {
	l.handlers = make([]Handler, 0)
}

func (l *logger) SetErrorHandler(w ErrorHandler) {
	l.errorHandler = w
}

func (l *logger) EnableCaller(level Level) {
	if level >= _MIN && level <= _MAX {
		l.caller[level] = true
	}
}

func (l *logger) DisableCaller(level Level) {
	if level >= _MIN && level <= _MAX {
		l.caller[level] = false
	}
}

func (l *logger) EnableStack(level Level) {
	if level >= _MIN && level <= _MAX {
		l.stack[level] = true
	}
}

func (l *logger) DisableStack(level Level) {
	if level >= _MIN && level <= _MAX {
		l.stack[level] = false
	}
}

func (l *logger) SetTag(k interface{}, v interface{}) {
	l.tags[k] = v
}

func (l *logger) DelTag(k interface{}) {
	delete(l.tags, k)
}

func (l *logger) ResetTag() {
	l.tags = make(map[interface{}]interface{})
}

func (l *logger) Tags() map[interface{}]interface{} {
	return l.tags
}

func (l *logger) SetKv(k interface{}, v interface{}) {
	l.kvs[k] = v
}

func (l *logger) DelKv(k interface{}) {
	delete(l.kvs, k)
}

func (l *logger) ResetKv() {
	l.kvs = make(map[interface{}]interface{})
}

func (l *logger) Kvs() map[interface{}]interface{} {
	return l.kvs
}

func (l *logger) Debug() *Event {
	return l.log(DEBUG)
}

func (l *logger) Info() *Event {
	return l.log(INFO)
}

func (l *logger) Warn() *Event {
	return l.log(WARN)
}

func (l *logger) Error() *Event {
	return l.log(ERROR)
}

func (l *logger) Panic() *Event {
	return l.log(PANIC)
}

func (l *logger) Fatal() *Event {
	return l.log(FATAL)
}

func (l *logger) Flush() {
	for _, handler := range l.handlers {
		if err := handler.Flush(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}

	// ignore error produced by errorHandler
	_ = l.errorHandler.Flush()
}

func (l *logger) Close() {
	for _, handler := range l.handlers {
		if err := handler.Close(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}

	// ignore error produced by errorHandler
	_ = l.errorHandler.Close()
}

func (l *logger) logCaller(level Level) bool {
	if need, ok := l.caller[level]; ok {
		return need
	}

	return false
}

func (l *logger) logStack(level Level) bool {
	if need, ok := l.stack[level]; ok {
		return need
	}

	return false
}

// couldEnd could end the Logger with panic or os.exit().
func (l *logger) couldEnd(level Level, v interface{}) {
	// Note: If there is any level bigger than PANIC added, the logic here should be updated.
	switch level {
	case PANIC:
		l.Flush()
		panic(v)
	case FATAL:
		l.Close()
		exit(1)
	}
}

func (l *logger) log(level Level) *Event {
	if level < l.level {
		l.couldEnd(level, "")
		// No need to generate an Event for and then be handled.
		return nil
	}

	return newEvent(l, level)
}

func (l *logger) handle(event *Event) {
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

package easylog

var m *manager
var root *logger

func init() {
	m = &manager{
		loggerMap: make(map[string]*logger),
	}

	root = m.getLogger("")

	m.root = root
}

func GetLogger(name string) *logger {
	return m.getLogger(name)
}

// GetRootLogger is equivalent to GetLogger("")
func GetRootLogger() *logger {
	return m.getLogger("")
}

func SetLevel(level Level) {
	root.SetLevel(level)
}

func GetLevel() Level {
	return root.GetLevel()
}

func EnableFrame(level Level) {
	root.EnableCaller(level)
}

func DisableFrame(level Level) {
	root.DisableCaller(level)
}

func AddFilter(f IFilter) {
	root.AddFilter(f)
}

func RemoveFilter(f IFilter) {
	root.RemoveFilter(f)
}

func AddHandler(h IHandler) {
	root.AddHandler(h)
}

func RemoveHandler(h IHandler) {
	root.RemoveHandler(h)
}

func Debug() *Event {
	return root.log(DEBUG)
}

func Info() *Event {
	return root.log(INFO)
}

func Warn() *Event {
	return root.log(WARN)
}

func Error() *Event {
	return root.log(ERROR)
}

func Fatal() *Event {
	return root.log(FATAL)
}

func Flush() {
	root.Flush()
}

func Close() {
	root.Close()
}

/*
func NewCachedLogger(parent *Logger) *Logger {
	return &Logger{
		name:        "",
		manager:     m,
		cached:      true,
		Level:       INFO,
		parent:      parent,
		propagate:   false,
		placeholder: false,
		children:    make(map[*Logger]interface{}),
	}
}

func SetCached(cached bool) {
	root.SetCached(cached)
}
*/

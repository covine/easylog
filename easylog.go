package easylog

var m *manager
var root *logger

func init() {
	m = &manager{
		loggerMap: make(map[string]*logger),
	}

	root = newLogger()
	root.name = "root"
	root.manager = m

	m.root = root
}

func GetLogger(name string) *logger {
	return m.getLogger(name)
}

func GetRootLogger() *logger {
	return m.root
}

func SetLevel(level Level) {
	root.SetLevel(level)
}

func GetLevel() Level {
	return root.GetLevel()
}

func EnableFrame(level Level) {
	root.EnableFrame(level)
}

func DisableFrame(level Level) {
	root.DisableFrame(level)
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
	return root.log(DEBUG, 2)
}

func Info() *Event {
	return root.log(INFO, 2)
}

func Warn() *Event {
	return root.log(WARN, 2)
}

func Error() *Event {
	return root.log(ERROR, 2)
}

func Fatal() *Event {
	return root.log(FATAL, 2)
}

func Flush() {
	root.Flush()
}

func Close() {
	root.Close()
}

/*
func NewCachedLogger(parent *logger) *logger {
	return &logger{
		name:        "",
		manager:     m,
		cached:      true,
		level:       INFO,
		parent:      parent,
		propagate:   false,
		placeholder: false,
		children:    make(map[*logger]interface{}),
	}
}

func SetCached(cached bool) {
	root.SetCached(cached)
}
*/

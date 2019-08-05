package easylog

var m *manager
var root *Logger

func init() {
	m = &manager{
		loggerMap: make(map[string]*Logger),
	}

	root = &Logger{
		name:           "root",
		manager:        m,
		level:          NOTSET,
		parent:         nil,
		propagate:      true,
		isPlaceholder:  false,
		placeholderMap: make(map[*Logger]interface{}),
	}
	m.root = root
}

func GetRootLogger() *Logger {
	return m.root
}

func GetLogger(name string) *Logger {
	return m.getLogger(name)
}

func NewCachedLogger(parent *Logger) *Logger {
	return &Logger{
		name:           "",
		manager:        m,
		cached:         true,
		level:          NOTSET,
		parent:         parent,
		propagate:      false,
		isPlaceholder:  false,
		placeholderMap: make(map[*Logger]interface{}),
	}
}

func SetLevel(level Level) {
	root.SetLevel(level)
}

func SetLevelByString(level string) {
	root.SetLevelByString(level)
}

func EnableFrame(level Level) {
	root.EnableFrame(level)
}

func DisableFrame(level Level) {
	root.DisableFrame(level)
}

func SetCached(cached bool) {
	root.SetCached(cached)
}

func AddFilter(fi IFilter) {
	root.AddFilter(fi)
}

func RemoveFilter(fi IFilter) {
	root.RemoveFilter(fi)
}

func AddHandler(hw IEasyLogHandler) {
	root.AddHandler(hw)
}

func RemoveHandler(hw IEasyLogHandler) {
	root.RemoveHandler(hw)
}

func Debug() *Record {
	return root.Debug()
}

func Info() *Record {
	return root.Info()
}

func Warning() *Record {
	return root.Warning()
}

func Warn() *Record {
	return root.Warn()
}

func Error() *Record {
	return root.Error()
}

func Fatal() *Record {
	return root.Fatal()
}

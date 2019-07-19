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

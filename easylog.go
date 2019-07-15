package easylog

var m *manager
var root *Logger

func init() {
	m = newManager()
	root = newRootLogger()
	root.setManager(m)
	m.setRoot(root)
}

func GetRootLogger() *Logger {
	return m.root
}

func GetLogger(name string) *Logger {
	return m.getLogger(name)
}

func GetSparkLogger() *Logger {
	s := newSparkLogger()
	s.setManager(m)
	s.setParent(root)
	return s
}

func Disable(level Level) {
	m.setDisable(level)
}

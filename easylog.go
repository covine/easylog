package easylog

var m *manager

func init() {
	m = newManager()
	root := newRootLogger()
	root.setManager(m)
	m.setRoot(root)
}

func GetRootLogger() *Logger {
	return m.root
}

func GetLogger(name string) *Logger {
	return m.getLogger(name)
}

func Disable(level Level) {
	m.setDisable(level)
}

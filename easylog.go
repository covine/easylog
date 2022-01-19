package easylog

var m *manager
var root *Logger

func init() {
	m = &manager{
		loggerMap: make(map[string]*Logger),
	}

	root = m.getLogger("")

	m.root = root
}

func GetLogger(name string) *Logger {
	return m.getLogger(name)
}

// GetRootLogger is equivalent to GetLogger("")
func GetRootLogger() *Logger {
	return m.getLogger("")
}

func SetLevel(level Level) {
	root.SetLevel(level)
}

func GetLevel() Level {
	return root.GetLevel()
}

func AddHandler(h Handler) {
	root.AddHandler(h)
}

func RemoveHandler(h Handler) {
	root.RemoveHandler(h)
}

func ResetHandler() {
	root.ResetHandler()
}

func SetErrorHandler(h ErrorHandler) {
	root.SetErrorHandler(h)
}

func EnableCaller(level Level) {
	root.EnableCaller(level)
}

func DisableCaller(level Level) {
	root.DisableCaller(level)
}

func EnableStack(level Level) {
	root.EnableStack(level)
}

func DisableStack(level Level) {
	root.DisableStack(level)
}

func SetTag(k string, v interface{}) {
	root.SetTag(k, v)
}

func DelTag(k string) {
	root.DelTag(k)
}

func ResetTag() {
	root.ResetTag()
}

func Tags() map[interface{}]interface{} {
	return root.Tags()
}

func SetKv(k interface{}, v interface{}) {
	root.SetKv(k, v)
}

func DelKv(k interface{}) {
	root.DelKv(k)
}

func ResetKv() {
	root.ResetKv()
}

func Kvs() map[interface{}]interface{} {
	return root.Kvs()
}

func Debug() *Event {
	return root.log(DEBUG, nil)
}

func Info() *Event {
	return root.log(INFO, nil)
}

func Warn() *Event {
	return root.log(WARN, nil)
}

func Error() *Event {
	return root.log(ERROR, nil)
}

func Panic() *Event {
	return root.log(PANIC, func(v interface{}) {
		root.Flush()
		panic(v)
	})
}

func Fatal() *Event {
	return root.log(FATAL, func(v interface{}) {
		root.Flush()
		root.Close()
		exit(1)
	})
}

func Flush() {
	root.Flush()
}

func Close() {
	root.Close()
}

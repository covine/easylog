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

func AddHandler(h Handler) {
	root.AddHandler(h)
}

func RemoveHandler(h Handler) {
	root.RemoveHandler(h)
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

func Tags() map[interface{}]interface{} {
	return root.Tags()
}

func SetKv(k interface{}, v interface{}) {
	root.SetKv(k, v)
}

func DelKv(k interface{}) {
	root.DelKv(k)
}

func Kvs() map[interface{}]interface{} {
	return root.Kvs()
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

func Panic() *Event {
	return root.log(PANIC)
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

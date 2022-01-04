package easylog

type Level int8

const (
	DEBUG Level = iota - 1 // -1
	INFO                   // 0
	WARN                   // 1
	ERROR                  // 2
	FATAL                  // 3
)

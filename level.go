package easylog

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	FATAL
)

package easylog

type IFilter interface {
	Filter(Record) bool
}

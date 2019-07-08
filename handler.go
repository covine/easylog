package easylog

type IHandler interface {
	SetLevel(Level)
	GetLevel() Level
	SetFormatter(IFormatter)
	AddFilter(IFilter)
	RemoveFilter(IFilter)
	Handle(record Record)
	Flush()
	Close()
}

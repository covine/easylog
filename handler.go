package easylog

import (
	"container/list"
)

type Formatter func(record Record) string

type IHandler interface {
	Handle(Record)
	Flush()
	Close()
}

type IEasyLogHandler interface {
	IHandler
	IFilters
	SetLevel(Level)
	GetLevel() Level
}

type handlerWrapper struct {
	Filters
	handler IHandler
	level   Level
}

func NewEasyLogHandler(ih IHandler) IEasyLogHandler {
	return &handlerWrapper{
		handler: ih,
	}
}

func (h *handlerWrapper) SetLevel(level Level) {
	if IsLevel(level) {
		h.level = level
	}
}

func (h *handlerWrapper) GetLevel() Level {
	return h.level
}

func (h *handlerWrapper) Handle(record Record) {
	if record.Level < h.level {
		return
	}

	if !h.Filters.Filter(record) {
		return
	}

	if h.handler != nil {
		h.handler.Handle(record)
	}
}

func (h *handlerWrapper) Flush() {
	if h.handler != nil {
		h.handler.Flush()
	}
}

func (h *handlerWrapper) Close() {
	if h.handler != nil {
		h.handler.Close()
	}
}

// not thread safe, set handlers during init
type Handlers struct {
	handlers *list.List
}

func (h *Handlers) AddHandler(hw IEasyLogHandler) {
	if hw == nil {
		return
	}

	if h.handlers == nil {
		h.handlers = list.New()
	}

	find := false
	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler == hw {
			find = true
			break
		}
	}

	if find {
		return
	} else {
		h.handlers.PushBack(hw)
	}
}

func (h *Handlers) RemoveHandler(hw IEasyLogHandler) {
	if hw == nil {
		return
	}

	if h.handlers == nil {
		h.handlers = list.New()
	}

	var next *list.Element
	for ele := h.handlers.Front(); ele != nil; ele = next {
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler == hw {
			next = ele.Next()
			h.handlers.Remove(ele)
		}
	}
}

func (h *Handlers) Handle(record Record) {
	if h.handlers == nil {
		return
	}

	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler != nil {
			handler.Handle(record)
		}
	}
}

func (h *Handlers) Flush() {
	if h.handlers == nil {
		return
	}
	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler != nil {
			handler.Flush()
		}
	}
}

func (h *Handlers) Close() {
	if h.handlers == nil {
		return
	}

	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler != nil {
			handler.Close()
		}
	}
}

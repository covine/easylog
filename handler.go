package easylog

import (
	"container/list"
)

type Formatter func(record *Event) string

// IHandler mockery --name=IHandler --inpackage --case=underscore --filename=handler_mock.go --structname MockHandler
type IHandler interface {
	Handle(*Event)
	Flush()
	Close()
}

type IEasyLogHandler interface {
	IHandler
	IFilters
	SetLevel(Level)
	GetLevel() Level
}

// NewEasyLogHandler
// Note: Event will be recycled after being handled, make sure do not cache Event in the Handler.
func NewEasyLogHandler(ih IHandler) IEasyLogHandler {
	return &handlerWrapper{
		handler: ih,
	}
}

type handlerWrapper struct {
	Filters
	handler IHandler
	level   Level
}

func (h *handlerWrapper) SetLevel(level Level) {
	h.level = level
}

func (h *handlerWrapper) GetLevel() Level {
	return h.level
}

func (h *handlerWrapper) Handle(record *Event) {
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

// Handlers
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
		next = ele.Next()
		handler, ok := ele.Value.(IEasyLogHandler)
		if ok && handler == hw {
			h.handlers.Remove(ele)
		}
	}
}

func (h *Handlers) HasHandler() bool {
	if h.handlers == nil {
		return false
	} else {
		return h.handlers.Len() > 0
	}
}

func (h *Handlers) Handle(record *Event) {
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

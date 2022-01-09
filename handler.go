package easylog

type Formatter func(record *Event) string

// IHandler mockery --name=IHandler --inpackage --case=underscore --filename=handler_mock.go --structname MockHandler
type IHandler interface {
	Handle(*Event)
	Flush()
	Close()
}

type Core interface {
	LevelEnabler

	// With adds structured context to the Core.
	With([]Field) Core
	// Check determines whether the supplied Entry should be logged (using the
	// embedded LevelEnabler and possibly some extra logic). If the entry
	// should be logged, the Core adds itself to the CheckedEntry and returns
	// the result.
	//
	// Callers must use Check before calling Write.
	Check(Entry, *CheckedEntry) *CheckedEntry
	// Write serializes the Entry and any Fields supplied at the log site and
	// writes them to their destination.
	//
	// If called, Write should always log the Entry and Fields; it should not
	// replicate the logic of Check.
	Write(Entry, []Field) error
	// Sync flushes buffered logs (if any).
	Sync() error
}

/*
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
	Level   Level
}

func (h *handlerWrapper) SetLevel(Level Level) {
	h.Level = Level
}

func (h *handlerWrapper) GetLevel() Level {
	return h.Level
}

func (h *handlerWrapper) Handle(record *Event) {
	if record.Level < h.Level {
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
		handler, OK := ele.Value.(IEasyLogHandler)
		if OK && handler == hw {
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
		handler, OK := ele.Value.(IEasyLogHandler)
		if OK && handler == hw {
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
		handler, OK := ele.Value.(IEasyLogHandler)
		if OK && handler != nil {
			handler.Handle(record)
		}
	}
}

func (h *Handlers) Flush() {
	if h.handlers == nil {
		return
	}
	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, OK := ele.Value.(IEasyLogHandler)
		if OK && handler != nil {
			handler.Flush()
		}
	}
}

func (h *Handlers) Close() {
	if h.handlers == nil {
		return
	}
	for ele := h.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, OK := ele.Value.(IEasyLogHandler)
		if OK && handler != nil {
			handler.Close()
		}
	}
}
*/

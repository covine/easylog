package easylog

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func clear() {
	SetLevel(INFO)
	SetErrorHandler(&nopErrorHandler{})
	ResetTag()
	ResetKv()
	ResetHandler()
	DisableCaller(DEBUG)
	DisableCaller(INFO)
	DisableCaller(WARN)
	DisableCaller(ERROR)
	DisableCaller(PANIC)
	DisableCaller(FATAL)
	DisableStack(DEBUG)
	DisableStack(INFO)
	DisableStack(WARN)
	DisableStack(ERROR)
	DisableStack(PANIC)
	DisableStack(FATAL)
}

func TestGetRootLogger(t *testing.T) {
	defer clear()

	r := GetRootLogger()

	assert.NotNil(t, r.manager)
	assert.Equal(t, "", r.Name())
	assert.Nil(t, r.parent)
	assert.False(t, r.propagate)
	assert.False(t, r.placeholder)
	assert.NotNil(t, r.tags)
	assert.NotNil(t, r.kvs)
	assert.False(t, r.logCaller(DEBUG))
	assert.False(t, r.logCaller(INFO))
	assert.False(t, r.logCaller(WARN))
	assert.False(t, r.logCaller(ERROR))
	assert.False(t, r.logCaller(PANIC))
	assert.False(t, r.logCaller(FATAL))
	assert.Equal(t, INFO, r.level)
	assert.NotNil(t, r.children)
	assert.Equal(t, 0, len(r.handlers))
	assert.Equal(t, nopErrorHandler{}, *(r.errorHandler.(*nopErrorHandler)))
}

func TestSetGetLevel(t *testing.T) {
	defer clear()

	SetLevel(DEBUG)

	assert.Equal(t, DEBUG, GetLevel())
}

func TestEnableDisableFrame(t *testing.T) {
	defer clear()

	EnableCaller(DEBUG)
	EnableCaller(INFO)
	EnableCaller(WARN)
	EnableCaller(ERROR)
	EnableCaller(PANIC)
	EnableCaller(FATAL)
	EnableCaller(-100)
	EnableCaller(100)
	assert.True(t, root.logCaller(DEBUG))
	assert.True(t, root.logCaller(INFO))
	assert.True(t, root.logCaller(WARN))
	assert.True(t, root.logCaller(ERROR))
	assert.True(t, root.logCaller(PANIC))
	assert.True(t, root.logCaller(FATAL))
	assert.False(t, root.logCaller(-100))
	assert.False(t, root.logCaller(100))
	DisableCaller(DEBUG)
	DisableCaller(INFO)
	DisableCaller(WARN)
	DisableCaller(ERROR)
	DisableCaller(PANIC)
	DisableCaller(FATAL)
	DisableCaller(-100)
	DisableCaller(100)
	assert.False(t, root.logCaller(DEBUG))
	assert.False(t, root.logCaller(INFO))
	assert.False(t, root.logCaller(WARN))
	assert.False(t, root.logCaller(ERROR))
	assert.False(t, root.logCaller(PANIC))
	assert.False(t, root.logCaller(FATAL))
	assert.False(t, root.logCaller(-100))
	assert.False(t, root.logCaller(100))

	Flush()
	Close()
}

func TestEnableDisableStack(t *testing.T) {
	defer clear()

	assert.Equal(t, false, root.logStack(DEBUG))
	assert.Equal(t, false, root.logStack(INFO))
	assert.Equal(t, false, root.logStack(WARN))
	assert.Equal(t, false, root.logStack(ERROR))
	assert.Equal(t, false, root.logStack(PANIC))
	assert.Equal(t, false, root.logStack(FATAL))

	EnableStack(DEBUG)
	EnableStack(INFO)
	EnableStack(WARN)
	EnableStack(ERROR)
	EnableStack(PANIC)
	EnableStack(FATAL)

	assert.Equal(t, true, root.logStack(DEBUG))
	assert.Equal(t, true, root.logStack(INFO))
	assert.Equal(t, true, root.logStack(WARN))
	assert.Equal(t, true, root.logStack(ERROR))
	assert.Equal(t, true, root.logStack(PANIC))
	assert.Equal(t, true, root.logStack(FATAL))

	DisableStack(DEBUG)
	DisableStack(INFO)
	DisableStack(WARN)
	DisableStack(ERROR)
	DisableStack(PANIC)
	DisableStack(FATAL)

	assert.Equal(t, false, root.logStack(DEBUG))
	assert.Equal(t, false, root.logStack(INFO))
	assert.Equal(t, false, root.logStack(WARN))
	assert.Equal(t, false, root.logStack(ERROR))
	assert.Equal(t, false, root.logStack(PANIC))
	assert.Equal(t, false, root.logStack(FATAL))
}

func TestNopErrorHandler(t *testing.T) {
	defer clear()

	assert.Equal(t, nopErrorHandler{}, *(root.errorHandler.(*nopErrorHandler)))

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == root && e.level == INFO && e.msg == "1" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, errors.New("ph handle error"))
	h.On("Flush").Return(errors.New("h flush error"))
	h.On("Close").Return(func() error { return nil })
	AddHandler(h)

	Info().Logf("1")

	Flush()
	Close()

	RemoveHandler(h)
	h.AssertExpectations(t)
}

func TestSetErrorHandler(t *testing.T) {
	defer clear()

	assert.Equal(t, nopErrorHandler{}, *(root.errorHandler.(*nopErrorHandler)))

	h := &MockErrorHandler{}

	SetErrorHandler(h)

	assert.Equal(t, h, root.errorHandler)
}

func TestTag(t *testing.T) {
	defer clear()

	assert.Equal(t, 0, len(Tags()))

	SetTag("service", "easylog")
	assert.Equal(t, 1, len(Tags()))
	assert.Equal(t, "easylog", Tags()["service"])

	SetTag("host", "vm")
	assert.Equal(t, 2, len(Tags()))
	assert.Equal(t, "easylog", Tags()["service"])
	assert.Equal(t, "vm", Tags()["host"])

	DelTag("service")
	assert.Equal(t, 1, len(Tags()))
	assert.Equal(t, "vm", Tags()["host"])

	DelTag("host")
	assert.Equal(t, 0, len(Tags()))
}

func TestKv(t *testing.T) {
	defer clear()

	assert.Equal(t, 0, len(Kvs()))

	SetKv("service", "easylog")
	assert.Equal(t, 1, len(Kvs()))
	assert.Equal(t, "easylog", Kvs()["service"])

	SetKv("host", "vm")
	assert.Equal(t, 2, len(Kvs()))
	assert.Equal(t, "easylog", Kvs()["service"])
	assert.Equal(t, "vm", Kvs()["host"])

	DelKv("service")
	assert.Equal(t, 1, len(Kvs()))
	assert.Equal(t, "vm", Kvs()["host"])

	DelKv("host")
	assert.Equal(t, 0, len(Kvs()))
}

func TestAddRemoveHandler(t *testing.T) {
	defer clear()

	assert.Equal(t, 0, len(root.handlers))

	AddHandler(nil)
	assert.Equal(t, 0, len(root.handlers))
	RemoveHandler(nil)
	assert.Equal(t, 0, len(root.handlers))

	n := NewNopHandler()

	AddHandler(n)
	assert.Equal(t, 1, len(root.handlers))
	assert.Equal(t, n, root.handlers[0])

	AddHandler(NewNopHandler())
	assert.Equal(t, 1, len(root.handlers))
	assert.Equal(t, n, root.handlers[0])

	h := &MockHandler{}

	AddHandler(h)
	assert.Equal(t, 2, len(root.handlers))
	assert.Equal(t, n, root.handlers[0])
	assert.Equal(t, h, root.handlers[1])

	AddHandler(h)
	assert.Equal(t, 2, len(root.handlers))
	assert.Equal(t, n, root.handlers[0])
	assert.Equal(t, h, root.handlers[1])

	RemoveHandler(h)
	assert.Equal(t, 1, len(root.handlers))
	assert.Equal(t, n, root.handlers[0])

	RemoveHandler(n)
	assert.Equal(t, 0, len(root.handlers))
}

func TestLog(t *testing.T) {
	defer clear()

	SetLevel(ERROR)

	l := GetLogger("t")

	pe := &MockErrorHandler{}
	pe.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "ph handle error"
	})).Once().Return(nil)
	pe.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(func(error) error { return nil })
	pe.On("Flush").Return(func() error { return errors.New("e flush error") })
	pe.On("Close").Return(func() error { return nil })
	SetErrorHandler(pe)

	ph := &MockHandler{}
	ph.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == ERROR && e.msg == "4" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(
		func(*Event) bool { return true },
		func(*Event) error { return errors.New("ph handle error") },
	)
	ph.On("Flush").Return(func() error { return errors.New("h flush error") })
	ph.On("Close").Return(nil)
	AddHandler(ph)

	pn := &MockHandler{}
	pn.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == ERROR && e.msg == "4" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(false, nil)
	pn.On("Flush").Return(nil)
	pn.On("Close").Return(nil)
	AddHandler(pn)

	pl := &MockHandler{}
	pl.On("Flush").Return(nil)
	pl.On("Close").Return(nil)
	AddHandler(pl)

	l.SetLevel(DEBUG)
	l.SetPropagate(true)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h handle error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "m close error"
	})).Once().Return(nil)
	e.On("Flush").Return(errors.New("e flush error"))
	e.On("Close").Return(nil)
	l.SetErrorHandler(e)

	n := NewNopHandler()
	l.AddHandler(n)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == DEBUG && e.msg == "1" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == INFO && e.msg == "2" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == WARN && e.msg == "3" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == ERROR && e.msg == "4" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Flush").Return(errors.New("h flush error"))
	h.On("Close").Return(nil)
	l.AddHandler(h)

	m := &MockHandler{}
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == INFO && e.msg == "2" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == WARN && e.msg == "3" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == l && e.level == ERROR && e.msg == "4" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(true, nil)
	m.On("Flush").Return(nil)
	m.On("Close").Return(errors.New("m close error"))
	l.AddHandler(m)

	l.Debug().Logf("1")
	l.Info().Logf("2")
	l.Warn().Logf("3")
	l.Error().Logf("4")

	l.Flush()
	l.Close()

	Flush()
	Close()

	m.AssertExpectations(t)
	h.AssertExpectations(t)
	e.AssertExpectations(t)
	pe.AssertExpectations(t)
	ph.AssertExpectations(t)
	pn.AssertExpectations(t)
	pl.AssertExpectations(t)
}

func TestPanicWithLevelINFO(t *testing.T) {
	defer clear()

	SetLevel(INFO)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h handle error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	e.On("Flush").Return(errors.New("e flush error"))
	SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == root && e.level == PANIC && e.msg == "1" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Flush").Return(errors.New("h flush error"))
	AddHandler(h)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "1", r.(string))
			h.AssertExpectations(t)
			e.AssertExpectations(t)
			return
		}
		assert.Fail(t, "should panic")
		return
	}()

	Panic().Logf("1")
	Panic().Logf("2")
}

func TestPanicWithLevelFATAL(t *testing.T) {
	defer clear()

	SetLevel(FATAL)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	e.On("Flush").Return(errors.New("e flush error"))
	SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Flush").Return(errors.New("h flush error"))
	AddHandler(h)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "", r.(string))
			h.AssertExpectations(t)
			e.AssertExpectations(t)
			return
		}
		assert.Fail(t, "should panic")
		return
	}()

	Panic().Logf("1")
	Panic().Logf("2")
}

func TestFatalWithLevelINFO(t *testing.T) {
	defer clear()

	f := &fakeExit{}
	modifyExit(f)
	defer recoverExit()

	SetLevel(INFO)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h handle error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h close error"
	})).Once().Return(nil)
	e.On("Flush").Return(nil)
	e.On("Close").Return(errors.New("e close error"))
	SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.GetLogger() == root && e.level == FATAL && e.msg == "1" &&
			e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false &&
			e.time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Flush").Return(nil)
	h.On("Close").Return(errors.New("h close error"))
	AddHandler(h)

	Fatal().Logf("%d", 1)

	assert.Equal(t, 1, f.code())
	h.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestFatalWithHighThanINFOLevel(t *testing.T) {
	defer clear()

	f := &fakeExit{}
	modifyExit(f)
	defer recoverExit()

	SetLevel(100)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h close error"
	})).Once().Return(nil)
	e.On("Flush").Return(nil)
	e.On("Close").Return(errors.New("e close error"))
	SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Close").Return(errors.New("h close error"))
	h.On("Flush").Return(nil)
	AddHandler(h)

	Fatal().Logf("1")

	assert.Equal(t, 1, f.code())
	h.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestDebug(t *testing.T) {
	defer clear()

	t.Run("emit Debug log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with DEBUG level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with WARN level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with invalid low level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with invalid high level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestInfo(t *testing.T) {
	defer clear()

	t.Run("emit Info log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with DEBUG level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with WARN level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with invalid low level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with invalid high level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestWarn(t *testing.T) {
	defer clear()

	t.Run("emit Warn log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == WARN &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Warn().Log()
		Warn().Log()
		Warn().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestError(t *testing.T) {
	defer clear()

	t.Run("emit Error log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == ERROR &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 && e.caller.ok == false
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Error().Log()
		Error().Log()
		Error().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Error log and enable frame", func(t *testing.T) {
		SetLevel(INFO)
		EnableCaller(ERROR)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == root && e.level == ERROR &&
				e.msg == "" &&
				e.extra == nil && e.caller.pc > 0 && strings.Contains(e.caller.file, "easylog_test.go") &&
				e.caller.ok == true
		})).Return(true, nil)
		m.On("Flush").Return(nil)
		m.On("Close").Return(nil)

		AddHandler(m)

		Error().Log()
		Error().Log()
		Error().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 1)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestGetLogger(t *testing.T) {
	defer clear()

	t.Run("get Logger", func(t *testing.T) {
		SetLevel(DEBUG)

		r := GetRootLogger()
		assert.NotNil(t, r.manager)
		assert.Equal(t, "", r.name)
		assert.Nil(t, r.parent)
		assert.False(t, r.propagate)
		assert.False(t, r.placeholder)
		assert.NotNil(t, r.tags)
		assert.NotNil(t, r.kvs)
		assert.Equal(t, DEBUG, r.level)
		assert.NotNil(t, r.children)
		assert.NotNil(t, r.handlers)
		assert.Equal(t, 0, len(r.children))

		empty := GetLogger("")
		assert.True(t, r == empty)
		assert.Equal(t, false, empty.propagate)
		assert.Equal(t, 0, len(empty.children))
		assert.Equal(t, false, empty.placeholder)
		assert.Equal(t, 0, len(r.children))

		emptyEmpty := GetLogger(".")
		assert.True(t, r == emptyEmpty.parent)
		assert.Equal(t, false, emptyEmpty.propagate)
		assert.Equal(t, 0, len(emptyEmpty.children))
		assert.Equal(t, false, emptyEmpty.placeholder)
		assert.Equal(t, 0, len(r.children))

		emptyA := GetLogger(".a")
		assert.True(t, root == emptyA.parent)
		assert.Equal(t, false, emptyA.propagate)
		assert.Equal(t, 0, len(emptyA.children))
		assert.Equal(t, false, emptyA.placeholder)
		assert.Equal(t, 0, len(r.children))

		emptyEmptyA := GetLogger("..a")
		assert.True(t, emptyEmpty == emptyEmptyA.parent)
		assert.Equal(t, false, emptyEmptyA.propagate)
		assert.Equal(t, 0, len(emptyEmptyA.children))
		assert.Equal(t, false, emptyEmptyA.placeholder)
		assert.Equal(t, 0, len(emptyEmpty.children))
		assert.Equal(t, 0, len(r.children))

		emptyEmptyAEmptyEmpty := GetLogger("..a..")
		assert.True(t, emptyEmptyA == emptyEmptyAEmptyEmpty.parent)
		assert.True(t, GetLogger("..a.") == emptyEmptyAEmptyEmpty.parent)
		assert.True(t, GetLogger("..a.").parent == emptyEmptyA)
		assert.Equal(t, false, GetLogger("..a.").placeholder)
		assert.Equal(t, 0, len(GetLogger("..a.").children))
		assert.Equal(t, 0, len(emptyEmptyA.children))
		assert.Equal(t, 0, len(emptyEmpty.children))
		assert.Equal(t, 0, len(r.children))

		a5 := GetLogger("a.b.c.d.e")
		assert.True(t, r == a5.parent)
		assert.True(t, GetLogger("a") == a5.parent)
		assert.Equal(t, root, GetLogger("a").parent)
		assert.Equal(t, false, GetLogger("a").placeholder)
		assert.Equal(t, 0, len(GetLogger("a").children))
		assert.Equal(t, 0, len(r.children))

		ab := GetLogger("a.b")
		assert.True(t, ab == a5.parent)
		assert.Equal(t, GetLogger("a"), GetLogger("a.b").parent)
		assert.Equal(t, 0, len(GetLogger("a").children))
		assert.Equal(t, 0, len(GetLogger("a.b").children))
		assert.Equal(t, 0, len(r.children))

		a4 := GetLogger("a.b.c.d")
		assert.True(t, a4.parent == ab)

		assert.True(t, a5.parent == a4)

		a7 := GetLogger("a.b.c.d.e.d.c")
		assert.True(t, a7.parent == a5)

		b7 := GetLogger("b.b.c.d.e.d.c")
		assert.True(t, b7.parent == root)

		c5 := GetLogger("1.2.3.4.5")
		assert.True(t, c5.parent == root)

		c4 := GetLogger("1.2.3.4")
		assert.True(t, c4.parent == root)
		assert.True(t, c5.parent == c4)

		c1 := GetLogger("1")
		c3 := GetLogger("1.2")
		assert.True(t, c3.parent == c1)
		assert.True(t, c4.parent == c3)

		fakeRoot := GetLogger("root")
		assert.True(t, r != fakeRoot)
		assert.True(t, r == fakeRoot.parent)
	})

	t.Run("get Logger concurrently", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10000; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()
				l := GetLogger(strconv.Itoa(j))
				assert.True(t, l.parent == GetRootLogger())
			}(i)
		}
		w.Wait()
	})
}

func TestLogger(t *testing.T) {
	defer clear()

	a := GetLogger("a")
	ab := GetLogger("a.b")
	abc := GetLogger("a.b.c")
	abcd := GetLogger("a.b.c.d")
	assert.True(t, root == a.parent)
	assert.True(t, a == ab.parent)
	assert.True(t, ab == abc.parent)
	assert.True(t, abc == abcd.parent)
	assert.True(t, nil == root.parent)

	assert.Equal(t, "a", a.Name())
	assert.Equal(t, false, a.GetPropagate())
	a.SetPropagate(true)
	assert.Equal(t, true, a.GetPropagate())

	assert.Equal(t, false, ab.GetPropagate())
	ab.SetPropagate(true)
	assert.Equal(t, true, ab.GetPropagate())

	assert.Equal(t, false, abc.GetPropagate())

	assert.Equal(t, false, abcd.GetPropagate())
	abcd.SetPropagate(true)
	assert.Equal(t, true, abcd.GetPropagate())

	assert.Equal(t, false, root.propagate)

	t.Run("event -> <a.b.c.d> -> <a.b.c>", func(t *testing.T) {
		m1 := &MockHandler{}
		m1.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == INFO &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m1.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == ERROR &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m1.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == WARN &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m1.On("Flush").Return(nil)
		m1.On("Close").Return(nil)
		abcd.AddHandler(m1)

		m2 := &MockHandler{}
		m2.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == INFO &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m2.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == ERROR &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m2.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.level == WARN &&
				e.caller.pc == 0 && e.caller.file == "" && e.caller.line == 0 &&
				e.caller.ok == false && e.time.Second() > 0
		})).Once().Return(true, nil)
		m2.On("Flush").Return(nil)
		m2.On("Close").Return(nil)
		abc.AddHandler(m2)

		m3 := &MockHandler{}
		ab.AddHandler(m3)

		m4 := &MockHandler{}
		a.AddHandler(m4)

		m5 := &MockHandler{}
		root.AddHandler(m5)

		abcd.Debug().Log()
		abcd.Info().Log()
		abcd.Error().Log()
		abcd.Warn().Log()

		root.RemoveHandler(m5)
		root.Flush()
		root.Close()

		a.RemoveHandler(m4)
		a.Flush()
		a.Close()

		ab.RemoveHandler(m3)
		ab.Flush()
		ab.Close()

		abc.RemoveHandler(m2)
		abc.Flush()
		abc.Close()

		abcd.RemoveHandler(m1)
		abcd.Flush()
		abcd.Close()
	})

	t.Run("event -> <a.b.c>", func(t *testing.T) {

	})

	t.Run("event -> <a.b> -> <a> -> <>", func(t *testing.T) {

	})

	t.Run("event -> <a> -> <>", func(t *testing.T) {

	})
}

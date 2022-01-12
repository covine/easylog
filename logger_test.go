package easylog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogger_Name(t *testing.T) {
	l := newLogger()
	assert.Equal(t, "", l.Name())
}

func TestLogger_SetGetPropagate(t *testing.T) {
	l := newLogger()

	assert.Equal(t, false, l.GetPropagate())

	l.SetPropagate(true)

	assert.Equal(t, true, l.GetPropagate())
}

func TestLogger_SetGetLevel(t *testing.T) {
	l := newLogger()

	assert.Equal(t, INFO, l.GetLevel())

	l.SetLevel(PANIC)

	assert.Equal(t, PANIC, l.GetLevel())
}

func TestLogger_EnableDisableCaller(t *testing.T) {
	l := newLogger()

	assert.Equal(t, false, l.logCaller(DEBUG))
	assert.Equal(t, false, l.logCaller(INFO))
	assert.Equal(t, false, l.logCaller(WARN))
	assert.Equal(t, false, l.logCaller(ERROR))
	assert.Equal(t, false, l.logCaller(PANIC))
	assert.Equal(t, false, l.logCaller(FATAL))

	l.EnableCaller(DEBUG)
	l.EnableCaller(INFO)
	l.EnableCaller(WARN)
	l.EnableCaller(ERROR)
	l.EnableCaller(PANIC)
	l.EnableCaller(FATAL)

	assert.Equal(t, true, l.logCaller(DEBUG))
	assert.Equal(t, true, l.logCaller(INFO))
	assert.Equal(t, true, l.logCaller(WARN))
	assert.Equal(t, true, l.logCaller(ERROR))
	assert.Equal(t, true, l.logCaller(PANIC))
	assert.Equal(t, true, l.logCaller(FATAL))

	l.DisableCaller(DEBUG)
	l.DisableCaller(INFO)
	l.DisableCaller(WARN)
	l.DisableCaller(ERROR)
	l.DisableCaller(PANIC)
	l.DisableCaller(FATAL)

	assert.Equal(t, false, l.logCaller(DEBUG))
	assert.Equal(t, false, l.logCaller(INFO))
	assert.Equal(t, false, l.logCaller(WARN))
	assert.Equal(t, false, l.logCaller(ERROR))
	assert.Equal(t, false, l.logCaller(PANIC))
	assert.Equal(t, false, l.logCaller(FATAL))
}

func TestLogger_EnableDisableStack(t *testing.T) {
	l := newLogger()

	assert.Equal(t, false, l.logStack(DEBUG))
	assert.Equal(t, false, l.logStack(INFO))
	assert.Equal(t, false, l.logStack(WARN))
	assert.Equal(t, false, l.logStack(ERROR))
	assert.Equal(t, false, l.logStack(PANIC))
	assert.Equal(t, false, l.logStack(FATAL))

	l.EnableStack(DEBUG)
	l.EnableStack(INFO)
	l.EnableStack(WARN)
	l.EnableStack(ERROR)
	l.EnableStack(PANIC)
	l.EnableStack(FATAL)

	assert.Equal(t, true, l.logStack(DEBUG))
	assert.Equal(t, true, l.logStack(INFO))
	assert.Equal(t, true, l.logStack(WARN))
	assert.Equal(t, true, l.logStack(ERROR))
	assert.Equal(t, true, l.logStack(PANIC))
	assert.Equal(t, true, l.logStack(FATAL))

	l.DisableStack(DEBUG)
	l.DisableStack(INFO)
	l.DisableStack(WARN)
	l.DisableStack(ERROR)
	l.DisableStack(PANIC)
	l.DisableStack(FATAL)

	assert.Equal(t, false, l.logStack(DEBUG))
	assert.Equal(t, false, l.logStack(INFO))
	assert.Equal(t, false, l.logStack(WARN))
	assert.Equal(t, false, l.logStack(ERROR))
	assert.Equal(t, false, l.logStack(PANIC))
	assert.Equal(t, false, l.logStack(FATAL))
}

func TestLogger_SetErrorHandler(t *testing.T) {
	l := newLogger()
	assert.Equal(t, nopErrorHandler{}, *(l.errorHandler.(*nopErrorHandler)))

	h := &MockErrorHandler{}

	l.SetErrorHandler(h)

	assert.Equal(t, h, l.errorHandler)
}

func TestLogger_Tag(t *testing.T) {
	l := newLogger()
	assert.Equal(t, 0, len(l.Tags()))

	l.SetTag("service", "easylog")
	assert.Equal(t, 1, len(l.Tags()))
	assert.Equal(t, "easylog", l.Tags()["service"])

	l.SetTag("host", "vm")
	assert.Equal(t, 2, len(l.Tags()))
	assert.Equal(t, "easylog", l.Tags()["service"])
	assert.Equal(t, "vm", l.Tags()["host"])

	l.DelTag("service")
	assert.Equal(t, 1, len(l.Tags()))
	assert.Equal(t, "vm", l.Tags()["host"])

	l.DelTag("host")
	assert.Equal(t, 0, len(l.Tags()))
}

func TestLogger_Kv(t *testing.T) {
	l := newLogger()
	assert.Equal(t, 0, len(l.Kvs()))

	l.SetKv("service", "easylog")
	assert.Equal(t, 1, len(l.Kvs()))
	assert.Equal(t, "easylog", l.Kvs()["service"])

	l.SetKv("host", "vm")
	assert.Equal(t, 2, len(l.Kvs()))
	assert.Equal(t, "easylog", l.Kvs()["service"])
	assert.Equal(t, "vm", l.Kvs()["host"])

	l.DelKv("service")
	assert.Equal(t, 1, len(l.Kvs()))
	assert.Equal(t, "vm", l.Kvs()["host"])

	l.DelKv("host")
	assert.Equal(t, 0, len(l.Kvs()))
}

func TestLogger_AddRemoveHandler(t *testing.T) {
	l := newLogger()

	assert.Equal(t, 0, len(l.handlers))

	l.AddHandler(nil)
	assert.Equal(t, 0, len(l.handlers))
	l.RemoveHandler(nil)
	assert.Equal(t, 0, len(l.handlers))

	n := &nopHandler{}

	l.AddHandler(n)
	assert.Equal(t, 1, len(l.handlers))
	assert.Equal(t, n, l.handlers[0])

	l.AddHandler(&nopHandler{})
	assert.Equal(t, 1, len(l.handlers))
	assert.Equal(t, n, l.handlers[0])

	h := &MockHandler{}

	l.AddHandler(h)
	assert.Equal(t, 2, len(l.handlers))
	assert.Equal(t, n, l.handlers[0])
	assert.Equal(t, h, l.handlers[1])

	l.AddHandler(h)
	assert.Equal(t, 2, len(l.handlers))
	assert.Equal(t, n, l.handlers[0])
	assert.Equal(t, h, l.handlers[1])

	l.RemoveHandler(h)
	assert.Equal(t, 1, len(l.handlers))
	assert.Equal(t, n, l.handlers[0])

	l.RemoveHandler(n)
	assert.Equal(t, 0, len(l.handlers))
}

func TestLogger_Log(t *testing.T) {
	p := newLogger()
	p.parent = nil
	p.SetLevel(ERROR)

	l := newLogger()
	l.parent = p

	pe := &MockErrorHandler{}
	pe.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "ph handle error"
	})).Once().Return(nil)
	pe.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	pe.On("Flush").Return(errors.New("e flush error"))
	pe.On("Close").Return(nil)
	p.SetErrorHandler(pe)

	ph := &MockHandler{}
	ph.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == ERROR && e.Msg == "4" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, errors.New("ph handle error"))
	ph.On("Flush").Return(errors.New("h flush error"))
	ph.On("Close").Return(nil)
	p.AddHandler(ph)

	pn := &MockHandler{}
	pn.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == ERROR && e.Msg == "4" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(false, nil)
	pn.On("Flush").Return(nil)
	pn.On("Close").Return(nil)
	p.AddHandler(pn)

	pl := &MockHandler{}
	pl.On("Flush").Return(nil)
	pl.On("Close").Return(nil)
	p.AddHandler(pl)

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

	n := &nopHandler{}
	l.AddHandler(n)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == DEBUG && e.Msg == "1" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == INFO && e.Msg == "2" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == WARN && e.Msg == "3" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == ERROR && e.Msg == "4" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, nil)
	h.On("Flush").Return(errors.New("h flush error"))
	h.On("Close").Return(nil)
	l.AddHandler(h)

	m := &MockHandler{}
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == INFO && e.Msg == "2" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, nil)
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == WARN && e.Msg == "3" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(true, nil)
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == ERROR && e.Msg == "4" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
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

	p.Flush()
	p.Close()

	m.AssertExpectations(t)
	h.AssertExpectations(t)
	e.AssertExpectations(t)
	pe.AssertExpectations(t)
	ph.AssertExpectations(t)
	pn.AssertExpectations(t)
	pl.AssertExpectations(t)
}

func TestLogger_PanicWithLevelINFO(t *testing.T) {
	l := newLogger()

	l.SetLevel(INFO)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h handle error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	e.On("Flush").Return(errors.New("e flush error"))
	l.SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == PANIC && e.Msg == "1" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Flush").Return(errors.New("h flush error"))
	l.AddHandler(h)

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

	l.Panic().Logf("1")
	l.Panic().Logf("2")
}

func TestLogger_PanicWithLevelFATAL(t *testing.T) {
	l := newLogger()

	l.SetLevel(FATAL)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h flush error"
	})).Once().Return(nil)
	e.On("Flush").Return(errors.New("e flush error"))
	l.SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Flush").Return(errors.New("h flush error"))
	l.AddHandler(h)

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

	l.Panic().Logf("1")
	l.Panic().Logf("2")
}

func TestLogger_FatalWithLevelINFO(t *testing.T) {
	f := &fakeExit{}
	modifyExit(f)
	defer recoverExit()

	l := newLogger()

	l.SetLevel(INFO)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h handle error"
	})).Once().Return(nil)
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h close error"
	})).Once().Return(nil)
	e.On("Close").Return(errors.New("e close error"))
	l.SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && e.Level == FATAL && e.Msg == "1" &&
			e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false &&
			e.Time.Second() > 0
	})).Once().Return(false, errors.New("h handle error"))
	h.On("Close").Return(errors.New("h close error"))
	l.AddHandler(h)

	l.Fatal().Logf("1")

	assert.Equal(t, 1, f.code())
	h.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestLogger_FatalWithHighThanINFOLevel(t *testing.T) {
	f := &fakeExit{}
	modifyExit(f)
	defer recoverExit()

	l := newLogger()

	l.SetLevel(100)

	e := &MockErrorHandler{}
	e.On("Handle", mock.MatchedBy(func(err error) bool {
		return err.Error() == "h close error"
	})).Once().Return(nil)
	e.On("Close").Return(errors.New("e close error"))
	l.SetErrorHandler(e)

	h := &MockHandler{}
	h.On("Close").Return(errors.New("h close error"))
	l.AddHandler(h)

	l.Fatal().Logf("1")

	assert.Equal(t, 1, f.code())
	h.AssertExpectations(t)
	e.AssertExpectations(t)
}

/*



func TestLogger_Caller(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)

	e := l.Info()
	e.Log()

	assert.True(t, e.OK)
	assert.True(t, e.PC > 0)

	l.DisableCaller(INFO)
}
*/

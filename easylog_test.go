package easylog

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRootLogger(t *testing.T) {
	r := GetRootLogger()

	assert.NotNil(t, r.manager)
	assert.Equal(t, "", r.Name())
	assert.Nil(t, r.parent)
	assert.False(t, r.propagate)
	assert.False(t, r.placeholder)
	assert.NotNil(t, r.tags)
	assert.NotNil(t, r.kvs)
	assert.False(t, r.debugCaller)
	assert.False(t, r.infoCaller)
	assert.False(t, r.warnCaller)
	assert.False(t, r.errorCaller)
	assert.False(t, r.fatalCaller)
	assert.Equal(t, INFO, r.level)
	assert.NotNil(t, r.children)
	assert.Equal(t, 0, len(*(r.filters)))
	assert.Equal(t, 0, len(*(r.handlers)))
}

func TestSetGetLevel(t *testing.T) {
	SetLevel(DEBUG)

	assert.Equal(t, DEBUG, GetLevel())
}

func TestEnableDisableFrame(t *testing.T) {
	EnableFrame(DEBUG)
	EnableFrame(INFO)
	EnableFrame(WARN)
	EnableFrame(ERROR)
	EnableFrame(FATAL)
	EnableFrame(-100)
	EnableFrame(100)
	assert.True(t, root.needCaller(DEBUG))
	assert.True(t, root.needCaller(INFO))
	assert.True(t, root.needCaller(WARN))
	assert.True(t, root.needCaller(ERROR))
	assert.True(t, root.needCaller(FATAL))
	assert.False(t, root.needCaller(-100))
	assert.False(t, root.needCaller(100))
	DisableFrame(DEBUG)
	DisableFrame(INFO)
	DisableFrame(WARN)
	DisableFrame(ERROR)
	DisableFrame(FATAL)
	DisableFrame(-100)
	DisableFrame(100)
	assert.False(t, root.needCaller(DEBUG))
	assert.False(t, root.needCaller(INFO))
	assert.False(t, root.needCaller(WARN))
	assert.False(t, root.needCaller(ERROR))
	assert.False(t, root.needCaller(FATAL))
	assert.False(t, root.needCaller(-100))
	assert.False(t, root.needCaller(100))

	Flush()
	Close()
}

func TestAddRemoveFilter(t *testing.T) {
	m1 := &MockFilter{}
	m2 := &MockFilter{}

	AddFilter(m1)
	AddFilter(m1)
	AddFilter(m1)
	AddFilter(m2)
	AddFilter(m2)
	AddFilter(m2)

	f := (*(root.filters))[0]
	assert.Equal(t, m1, f.(IFilter))
	assert.Equal(t, 2, len(*(root.filters)))
	f = (*(root.filters))[1]
	assert.Equal(t, m2, f.(IFilter))
	assert.Equal(t, 2, len(*(root.filters)))

	RemoveFilter(m1)
	f = (*(root.filters))[0]
	assert.Equal(t, m2, f.(IFilter))
	assert.Equal(t, 1, len(*(root.filters)))

	RemoveFilter(m2)
	assert.Equal(t, 0, len(*(root.filters)))

	Flush()
	Close()
}

func TestAddRemoveHandler(t *testing.T) {
	m1 := &MockHandler{}
	m2 := &MockHandler{}

	AddHandler(m1)
	AddHandler(m1)
	AddHandler(m1)
	AddHandler(m2)
	AddHandler(m2)
	AddHandler(m2)

	h := (*(root.handlers))[0]
	assert.Equal(t, m1, h.(IHandler))
	assert.Equal(t, 2, len(*(root.handlers)))
	h = (*(root.handlers))[1]
	assert.Equal(t, m2, h.(IHandler))
	assert.Equal(t, 2, len(*(root.handlers)))

	RemoveHandler(m1)
	h = (*(root.handlers))[0]
	assert.Equal(t, m2, h.(IHandler))
	assert.Equal(t, 1, len(*(root.handlers)))

	RemoveHandler(m2)
	assert.Equal(t, 0, len(*(root.handlers)))

	Flush()
	Close()
}

func TestDebug(t *testing.T) {
	t.Run("emit Debug log with INFO Level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with DEBUG Level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == DEBUG &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with WARN Level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == DEBUG &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with invalid low Level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == DEBUG &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Debug log with invalid high Level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == DEBUG &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Debug().Log()
		Debug().Log()
		Debug().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestInfo(t *testing.T) {
	t.Run("emit Info log with INFO Level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == INFO &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with DEBUG Level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == INFO &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with WARN Level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == INFO &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with invalid low Level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == INFO &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Info log with invalid high Level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == INFO &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Info().Log()
		Info().Log()
		Info().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 0)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestWarn(t *testing.T) {
	t.Run("emit Warn log with INFO Level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == WARN &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Warn().Log()
		Warn().Log()
		Warn().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestError(t *testing.T) {
	t.Run("emit Error log with INFO Level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == ERROR &&
				e.Msg == "" &&
				e.Extra == nil && e.PC == 0 && e.File == "" && e.Line == 0 && e.OK == false
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Error().Log()
		Error().Log()
		Error().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})

	t.Run("emit Error log and enable frame", func(t *testing.T) {
		SetLevel(INFO)
		EnableFrame(ERROR)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == ERROR &&
				e.Msg == "" &&
				e.Extra == nil && e.PC > 0 && strings.Contains(e.File, "easylog_test.go") &&
				e.Line >= 495 && e.Line <= 497 && e.OK == true
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		AddHandler(m)

		Error().Log()
		Error().Log()
		Error().Log()

		Flush()
		Close()

		RemoveHandler(m)

		m.AssertNumberOfCalls(t, "Handle", 3)
		m.AssertNumberOfCalls(t, "Flush", 2)
		m.AssertNumberOfCalls(t, "Close", 1)
	})
}

func TestFatal(t *testing.T) {
	t.Run("emit Fatal log and enable frame", func(t *testing.T) {
		SetLevel(INFO)
		EnableFrame(FATAL)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == root && e.Level == FATAL &&
				e.Msg == "" &&
				e.Extra == nil && e.PC > 0 && strings.Contains(e.File, "easylog_test.go") &&
				e.Line == 538 && e.OK == true
		})).Return()
		m.On("Flush").Return()
		m.On("Close").Return()

		defer func() {
			if r := recover(); r != nil {
				m.AssertNumberOfCalls(t, "Handle", 1)
				m.AssertNumberOfCalls(t, "Flush", 0)
				m.AssertNumberOfCalls(t, "Close", 0)
				RemoveHandler(m)
				return
			}
			assert.Fail(t, "should panic")
		}()

		AddHandler(m)

		Fatal().Log()
		Fatal().Log()
		Fatal().Log()

		Flush()
		Close()
	})
}

func TestGetLogger(t *testing.T) {
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
		assert.True(t, r.filters == nil || 0 == len(*(r.filters)))
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

func TestFilterHandlerParallel(t *testing.T) {
	m1 := &MockFilter{}
	m2 := &MockFilter{}

	h1 := &MockHandler{}
	h2 := &MockHandler{}

	l := GetLogger("test_filter_parallel")
	l.SetLevel(WARN)

	m1.On("Filter", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && (e.Level == WARN || e.Level == ERROR)
	})).Return(true)

	m2.On("Filter", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && (e.Level == WARN || e.Level == ERROR)
	})).Return(true)

	h1.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && (e.Level == WARN || e.Level == ERROR)
	})).Return()
	h1.On("Flush").Return()
	h1.On("Close").Return()

	h2.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.Logger == l && (e.Level == WARN || e.Level == ERROR)
	})).Return()
	h2.On("Flush").Return()
	h2.On("Close").Return()

	timeout := 7
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.AddFilter(m1)
				l.AddFilter(m2)
				l.AddFilter(nil)
				l.AddHandler(h1)
				l.AddHandler(h2)
				l.AddHandler(nil)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.RemoveFilter(m1)
				l.RemoveFilter(m2)
				l.RemoveFilter(nil)
				l.RemoveHandler(h1)
				l.RemoveHandler(h2)
				l.RemoveHandler(nil)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.Debug().Log()
				l.Info().Log()
				l.Warn().Log()
				l.Error().Log()
				l.Flush()
			}
		}
	}()

	wg.Wait()

	l.AddHandler(h1)
	l.AddHandler(h2)

	l.Flush()
	l.Close()

	m1.AssertExpectations(t)
	m2.AssertExpectations(t)
	h1.AssertExpectations(t)
	h2.AssertExpectations(t)

	l.RemoveFilter(m1)
	l.RemoveFilter(m2)
	l.RemoveHandler(h1)
	l.RemoveHandler(h2)
}

func TestLogger(t *testing.T) {
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
		f1 := &LevelEqualFilter{Level: ERROR}
		m1 := &MockHandler{}
		m1.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == abcd && e.Level == ERROR &&
				e.PC == 0 && e.File == "" && e.Line == 0 &&
				e.OK == false && e.Time.Second() > 0
		})).Once().Return()
		m1.On("Flush").Return()
		m1.On("Close").Return()
		abcd.AddFilter(f1)
		abcd.AddHandler(m1)

		f2 := &LevelEqualFilter{Level: ERROR}
		m2 := &MockHandler{}
		m2.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.Logger == abcd && e.Level == ERROR &&
				e.PC == 0 && e.File == "" && e.Line == 0 &&
				e.OK == false && e.Time.Second() > 0
		})).Once().Return()
		m2.On("Flush").Return()
		m2.On("Close").Return()
		abc.AddFilter(f2)
		abc.AddHandler(m2)

		f3 := &MockFilter{}
		m3 := &MockHandler{}
		ab.AddFilter(f3)
		ab.AddHandler(m3)

		f4 := &MockFilter{}
		m4 := &MockHandler{}
		a.AddFilter(f4)
		a.AddHandler(m4)

		f5 := &MockFilter{}
		m5 := &MockHandler{}
		root.AddFilter(f5)
		root.AddHandler(m5)
		//
		abcd.Debug().Log()
		abcd.Info().Log()
		abcd.Error().Log()
		abcd.Warn().Log()

		//
		root.RemoveFilter(f5)
		root.RemoveHandler(m5)
		root.Flush()
		root.Close()

		a.RemoveFilter(f4)
		a.RemoveHandler(m4)
		a.Flush()
		a.Close()

		ab.RemoveFilter(f3)
		ab.RemoveHandler(m3)
		ab.Flush()
		ab.Close()

		abc.RemoveFilter(f2)
		abc.RemoveHandler(m2)
		abc.Flush()
		abc.Close()

		abcd.RemoveFilter(f1)
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

func TestEvent(t *testing.T) {
	l := GetLogger("test_event")
	l.SetLevel(INFO)

	l.Debug().Log()
	l.Debug().Logf("debug")

	l.Info().Tag("k", "v").Log()
	l.Warn().Kv("a", "b").Logf("hi")
	l.Error().Attach("hi").Log()

	assert.NotPanics(t, func() {
		n := GetLogger("test_nil_event")
		n.SetLevel(ERROR)

		n.Debug().Tag("k", "v").Log()
		n.Info().Kv("k", "v").Logf("")
		n.Warn().Log()
		n.Warn().Logf("")
		n.Warn().Attach("").Log()
	})

	h := GetLogger("test_event_handler")
	h.SetLevel(INFO)
	h.EnableCaller(WARN)

	m := &MockHandler{}
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		t, ok := e.Tags["a"]
		if !ok {
			return false
		}

		k, ok := e.Kvs["c"]
		if !ok {
			return false
		}

		return e.Logger == h && e.Level == INFO &&
			e.Msg == "" && t == "b" && k == "d" &&
			e.Extra == "e" && e.PC == 0 &&
			e.File == "" && e.Line == 0 &&
			e.OK == false && e.Time.Second() > 0
	})).Once().Return()

	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		t, ok := e.Tags["a"]
		if !ok {
			return false
		}

		k, ok := e.Kvs["c"]
		if !ok {
			return false
		}

		return e.Logger == h && e.Level == WARN &&
			e.Msg == "f" && t == "b" && k == "d" &&
			e.Extra == "e" && e.PC > 0 &&
			strings.Contains(e.File, "event_test.go") &&
			e.Line == 81 && e.OK == true && e.Time.Second() > 0
	})).Once().Return()

	m.On("Flush").Return()
	m.On("Close").Return()

	h.AddHandler(m)

	h.Debug().Tag("a", "b").Kv("c", "d").Attach("e").Log()
	h.Info().Tag("a", "b").Kv("c", "d").Attach("e").Log()
	h.Warn().Tag("a", "b").Kv("c", "d").Attach("e").Logf("f")

	h.Flush()
	h.Close()
	h.RemoveHandler(m)
	h.DisableCaller(WARN)

	m.AssertExpectations(t)

	assert.Panics(t, func() {
		h.Fatal().Log()
	})
}

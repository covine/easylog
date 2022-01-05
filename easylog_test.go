package easylog

import (
	"strconv"
	"strings"
	"sync"
	"testing"

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
	assert.False(t, r.debugFrame)
	assert.False(t, r.infoFrame)
	assert.False(t, r.warnFrame)
	assert.False(t, r.errorFrame)
	assert.False(t, r.fatalFrame)
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
	assert.True(t, root.needFrame(DEBUG))
	assert.True(t, root.needFrame(INFO))
	assert.True(t, root.needFrame(WARN))
	assert.True(t, root.needFrame(ERROR))
	assert.True(t, root.needFrame(FATAL))
	assert.False(t, root.needFrame(-100))
	assert.False(t, root.needFrame(100))
	DisableFrame(DEBUG)
	DisableFrame(INFO)
	DisableFrame(WARN)
	DisableFrame(ERROR)
	DisableFrame(FATAL)
	DisableFrame(-100)
	DisableFrame(100)
	assert.False(t, root.needFrame(DEBUG))
	assert.False(t, root.needFrame(INFO))
	assert.False(t, root.needFrame(WARN))
	assert.False(t, root.needFrame(ERROR))
	assert.False(t, root.needFrame(FATAL))
	assert.False(t, root.needFrame(-100))
	assert.False(t, root.needFrame(100))

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
	t.Run("emit Debug log with INFO level", func(t *testing.T) {
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

	t.Run("emit Debug log with DEBUG level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Debug log with WARN level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Debug log with invalid low level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Debug log with invalid high level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == DEBUG &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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
	t.Run("emit Info log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Info log with DEBUG level", func(t *testing.T) {
		SetLevel(DEBUG)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Info log with WARN level", func(t *testing.T) {
		SetLevel(WARN)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Info log with invalid low level", func(t *testing.T) {
		SetLevel(-2)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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

	t.Run("emit Info log with invalid high level", func(t *testing.T) {
		SetLevel(10)

		m := &MockHandler{}

		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == INFO &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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
	t.Run("emit Warn log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == WARN &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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
	t.Run("emit Error log with INFO level", func(t *testing.T) {
		SetLevel(INFO)

		m := &MockHandler{}
		m.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.logger == root && e.level == ERROR &&
				e.msg == "" &&
				e.extra == nil && e.pc == 0 && e.file == "" && e.line == 0 && e.ok == false
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
			return e.logger == root && e.level == ERROR &&
				e.msg == "" &&
				e.extra == nil && e.pc > 0 && strings.Contains(e.file, "easylog_test.go") &&
				e.line >= 495 && e.line <= 497 && e.ok == true
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
			return e.logger == root && e.level == FATAL &&
				e.msg == "" &&
				e.extra == nil && e.pc > 0 && strings.Contains(e.file, "easylog_test.go") &&
				e.line == 538 && e.ok == true
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
	t.Run("get logger", func(t *testing.T) {
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

	t.Run("get logger concurrently", func(t *testing.T) {
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

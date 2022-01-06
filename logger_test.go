package easylog

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFilterHandlerParallel(t *testing.T) {
	m1 := &MockFilter{}
	m2 := &MockFilter{}

	h1 := &MockHandler{}
	h2 := &MockHandler{}

	l := GetLogger("test_filter_parallel")
	l.SetLevel(WARN)

	m1.On("Filter", mock.MatchedBy(func(e *Event) bool {
		return e.logger == l && (e.level == WARN || e.level == ERROR)
	})).Return(true)

	m2.On("Filter", mock.MatchedBy(func(e *Event) bool {
		return e.logger == l && (e.level == WARN || e.level == ERROR)
	})).Return(true)

	h1.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.logger == l && (e.level == WARN || e.level == ERROR)
	})).Return()
	h1.On("Flush").Return()
	h1.On("Close").Return()

	h2.On("Handle", mock.MatchedBy(func(e *Event) bool {
		return e.logger == l && (e.level == WARN || e.level == ERROR)
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
			return e.GetLogger() == abcd && e.GetLevel() == ERROR &&
				e.GetPC() == 0 && e.GetFile() == "" && e.GetLine() == 0 &&
				e.GetOK() == false && e.GetTime().Second() > 0
		})).Once().Return()
		m1.On("Flush").Return()
		m1.On("Close").Return()
		abcd.AddFilter(f1)
		abcd.AddHandler(m1)

		f2 := &LevelEqualFilter{Level: ERROR}
		m2 := &MockHandler{}
		m2.On("Handle", mock.MatchedBy(func(e *Event) bool {
			return e.GetLogger() == abcd && e.GetLevel() == ERROR &&
				e.GetPC() == 0 && e.GetFile() == "" && e.GetLine() == 0 &&
				e.GetOK() == false && e.GetTime().Second() > 0
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

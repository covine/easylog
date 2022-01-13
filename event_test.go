package easylog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilEvent(t *testing.T) {
	var e *Event = nil
	assert.NotPanics(t, func() {
		e.Tag("t", "v").Attach(nil).Log()
		e.Kv("k", "v").Logf("")
	})
}

func TestEventWithNilLogger(t *testing.T) {
	e := newEvent(nil, INFO)
	defer putEvent(e)

	assert.Equal(t, INFO, e.GetLevel())
	assert.Equal(t, nil, e.GetExtra())

	assert.NotPanics(t, func() {
		e.Tag("t", "v").Attach(nil)
		e.Kv("k", "v").Attach(nil)
		e.Attach(1)
	})

	assert.Equal(t, 1, e.GetExtra())
	assert.Equal(t, "", e.GetMsg())

	assert.Equal(t, "v", e.GetTags()["t"])
	assert.Equal(t, "v", e.GetKvs()["k"])

	assert.Panics(t, func() {
		e.Tag("t", "v").Attach(nil).Log()
	})

	assert.Panics(t, func() {
		e.Kv("k", "v").Attach(nil).Logf("")
	})
}

func TestEventWithCaller(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)
	e := newEvent(l, INFO)

	assert.Equal(t, l, e.GetLogger())

	e.Log()

	assert.True(t, e.GetCaller().GetOK())
	assert.True(t, e.GetCaller().GetPC() > 0)
	assert.True(t, strings.Contains(e.GetCaller().GetFile(), "event_test.go"))
	assert.True(t, e.GetCaller().GetFunc() == "github.com/covine/easylog.TestEventWithCaller")
}

func TestEventWithGetCallerError(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)
	e := newEvent(l, INFO)

	e.log("", 7000000)

	assert.False(t, e.caller.ok)
	assert.True(t, e.caller.pc == 0)
	assert.Equal(t, "", e.GetCaller().GetFile())
	assert.Equal(t, 0, e.GetCaller().GetLine())
	assert.Equal(t, "", e.GetCaller().GetFunc())
}

func TestEventWithStack(t *testing.T) {
	l := newLogger()
	l.EnableStack(INFO)

	e := newEvent(l, INFO)

	e.Log()

	assert.True(t, e.GetStack() != "")
	assert.True(t, e.GetTime().Second() > 0)
}

func deepStack(n int, e *Event) {
	if n <= 0 {
		e.Log()
	} else {
		deepStack(n-1, e)
	}
}

func TestEventWithDeepStack(t *testing.T) {
	l := newLogger()
	l.EnableStack(INFO)

	e := newEvent(l, INFO)

	deepStack(1000, e)

	assert.True(t, e.stack != "")
}

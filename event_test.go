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

	assert.NotPanics(t, func() {
		e.Tag("t", "v").Attach(nil)
		e.Kv("k", "v").Attach(nil)
	})

	assert.Panics(t, func() {
		e.Tag("t", "v").Attach(nil).Log()

	})

	assert.Panics(t, func() {
		e.Kv("k", "v").Attach(nil).Logf("")
	})

	putEvent(e)
}

func TestEventWithCaller(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)
	e := newEvent(l, INFO)

	e.Log()

	assert.True(t, e.OK)
	assert.True(t, e.PC > 0)
	assert.True(t, strings.Contains(e.File, "event_test.go"))
	assert.True(t, e.Func == "github.com/covine/easylog.TestEventWithCaller")
}

func TestEventWithGetCallerError(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)
	e := newEvent(l, INFO)

	e.log("", 7000000)

	assert.False(t, e.OK)
	assert.True(t, e.PC == 0)
	assert.Equal(t, "", e.File)
	assert.Equal(t, 0, e.Line)
	assert.Equal(t, "", e.Func)
}

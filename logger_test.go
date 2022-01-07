package easylog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Caller(t *testing.T) {
	l := newLogger()
	l.EnableCaller(INFO)

	e := l.Info()
	e.Log()

	assert.True(t, e.OK)
	assert.True(t, e.PC > 0)

	l.DisableCaller(INFO)
}

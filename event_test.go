package easylog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEvent(t *testing.T) {
	l := GetLogger("test_event")
	l.SetLevel(INFO)

	l.Debug().Log()
	l.Debug().Logf("debug")

	l.Info().Tag("k", "v").Log()
	l.Warn().Kv("a", "b").Logf("hi")
	l.Error().Extra("hi").Log()

	assert.NotPanics(t, func() {
		n := GetLogger("test_nil_event")
		n.SetLevel(ERROR)

		n.Debug().Tag("k", "v").Log()
		n.Info().Kv("k", "v").Logf("")
		n.Warn().Log()
		n.Warn().Logf("")
		n.Warn().Extra("").Log()
	})

	h := GetLogger("test_event_handler")
	h.SetLevel(INFO)
	h.EnableFrame(WARN)

	m := &MockHandler{}
	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		t, ok := e.GetTags().Load("a")
		if !ok {
			return false
		}

		k, ok := e.GetKvs().Load("c")
		if !ok {
			return false
		}

		return e.GetLogger() == h && e.GetLevel() == INFO &&
			e.GetMsg() == "" && t == "b" && k == "d" &&
			e.GetExtra() == "e" && e.GetPC() == 0 &&
			e.GetFile() == "" && e.GetLine() == 0 &&
			e.GetOK() == false && e.GetTime().Second() > 0
	})).Once().Return()

	m.On("Handle", mock.MatchedBy(func(e *Event) bool {
		t, ok := e.GetTags().Load("a")
		if !ok {
			return false
		}

		k, ok := e.GetKvs().Load("c")
		if !ok {
			return false
		}

		return e.GetLogger() == h && e.GetLevel() == WARN &&
			e.GetMsg() == "f" && t == "b" && k == "d" &&
			e.GetExtra() == "e" && e.GetPC() > 0 &&
			strings.Contains(e.GetFile(), "event_test.go") &&
			e.GetLine() == 81 && e.GetOK() == true && e.GetTime().Second() > 0
	})).Once().Return()

	m.On("Flush").Return()
	m.On("Close").Return()

	h.AddHandler(m)

	h.Debug().Tag("a", "b").Kv("c", "d").Extra("e").Log()
	h.Info().Tag("a", "b").Kv("c", "d").Extra("e").Log()
	h.Warn().Tag("a", "b").Kv("c", "d").Extra("e").Logf("f")

	h.Flush()
	h.Close()
	h.RemoveHandler(m)
	h.DisableFrame(WARN)

	m.AssertExpectations(t)

	assert.Panics(t, func() {
		h.Fatal().Log()
	})
}

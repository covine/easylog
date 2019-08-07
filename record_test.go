package easylog_test

import (
	"testing"

	"git.qutoutiao.net/govine/easylog"
)

func TestRecord_NIL(t *testing.T) {
	l := easylog.GetLogger("record.nil")
	l.SetLevel(easylog.WARN)
	r := l.Debug()
	if r != nil {
		t.Errorf("record lower than level must be nil")
		return
	}

	if l.Debug().ExistTag("hello") {
		t.Errorf("nil record do not exsit tag")
		return
	}

	l.Debug().Fields(map[string]interface{}{"hello": "world"}).Msg("hello world")
	l.Debug().Tag("hello").Msg("hello world")
	l.Debug().Extra("hello").Msg("hello world")
	l.Debug().Msg("hello world")
}

func TestRecord_ExistTag(t *testing.T) {
	l := easylog.GetLogger("record.ExistTag")
	l.SetLevel(easylog.DEBUG)
	r := l.Debug().Tag("hello")
	if !r.ExistTag("hello") {
		t.Errorf("hello tag must exists")
		return
	}
}

func TestRecord_Msg(t *testing.T) {
	l := easylog.GetLogger("record.Msg")
	l.SetLevel(easylog.WARN)
	l.Debug().Msg("hello")
	l.Debug().Msg("hello: %s", "world")
	l.Warn().Msg("hello")
	l.Warn().Msg("hello: %s", "world")
}

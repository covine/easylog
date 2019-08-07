package easylog_test

import (
	"testing"

	"github.com/govine/easylog"
	"github.com/govine/easylog/filter"
)

func TestLogger_SetLevelByString(t *testing.T) {
	l := easylog.GetLogger("logger.SetLevelByString")
	l.SetLevelByString("DEBUG")
	l.SetLevelByString("INFO")
	l.SetLevelByString("WARN")
	l.SetLevelByString("WARNING")
	l.SetLevelByString("ERROR")
	l.SetLevelByString("FATAL")
	l.SetLevelByString("")
}

func TestLogger_handlerRecord(t *testing.T) {
	l := easylog.GetLogger("logger.handlerRecord")
	l.SetLevelByString("WARN")
	l.AddFilter(&filter.LevelEqualFilter{Level: easylog.FATAL})

	c := easylog.GetLogger("logger.handlerRecord.c")
	c.SetLevel(easylog.DEBUG)
	c.SetPropagate(true)
	c.Debug().Msg("hello")
	c.Warn().Msg("hello")
}

func TestLogger_Flush(t *testing.T) {
	l := easylog.GetLogger("logger.Flush")
	l.SetLevelByString("WARN")
	l.SetCached(true)
	l.SetPropagate(false)
	l.Warn().Msg("hello")
	l.Close()
}

package easylog_test

import (
	"testing"

	"github.com/govine/easylog"
	"github.com/govine/easylog/filter"
)

func TestFilter(t *testing.T) {
	t.Run("filter", func(t *testing.T) {
		l := easylog.GetLogger("filter")
		f := &filter.LevelEqualFilter{Level: easylog.DEBUG}
		l.RemoveFilter(f)
		l.AddFilter(nil)
		l.AddFilter(f)
		l.AddFilter(f)
		l.RemoveFilter(nil)
		l.RemoveFilter(f)
		l.RemoveFilter(f)
	})
}

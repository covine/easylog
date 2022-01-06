package easylog

import (
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("filter", func(t *testing.T) {
		l := GetLogger("filter")
		f := &LevelEqualFilter{Level: DEBUG}
		l.RemoveFilter(f)
		l.AddFilter(nil)
		l.AddFilter(f)
		l.AddFilter(f)
		l.RemoveFilter(nil)
		l.RemoveFilter(f)
		l.RemoveFilter(f)
	})
}

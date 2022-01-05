package filter

import (
	"github.com/covine/easylog"
)

type LevelEqualFilter struct {
	Level easylog.Level
}

func (l *LevelEqualFilter) Filter(record *easylog.Event) bool {
	return record.GetLevel() == l.Level
}

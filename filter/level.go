package filter

import (
	"github.com/govine/easylog"
)

type LevelEqualFilter struct {
	Level easylog.Level
}

func (l *LevelEqualFilter) Filter(record *easylog.Record) bool {
	return record.Level == l.Level
}

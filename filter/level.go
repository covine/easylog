package filter

import (
	"git.qutoutiao.net/govine/easylog"
)

type LevelEqualFilter struct {
	Level easylog.Level
}

func (l *LevelEqualFilter) Filter(record easylog.Record) bool {
	return record.Level == l.Level
}

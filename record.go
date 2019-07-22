package easylog

import (
	"time"
)

type Record struct {
	Time  time.Time
	Level Level
	Msg   string
	Args  []interface{}
	File  string
	Line  int
}

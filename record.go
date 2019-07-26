package easylog

import (
	"sync"
	"time"
)

// methods of Record object are not concurrency-safe
type Record struct {
	Logger *Logger // which Logger Record belongs to

	Time      time.Time
	Level     Level
	Message   string
	Args      []interface{}
	FieldMap  map[string]interface{}
	Tags      []string
	ExtraData interface{}

	// runtime stack
	PC   uintptr
	File string
	Line int
	OK   bool
}

var recordPool = &sync.Pool{
	New: func() interface{} {
		return &Record{}
	},
}

func newRecord() *Record {
	r := recordPool.Get().(*Record)

	r.Logger = nil

	r.Time = time.Time{}
	r.Level = NOTSET
	r.Message = ""
	r.Args = nil
	r.FieldMap = nil
	r.Tags = nil
	r.ExtraData = nil

	r.PC = 0
	r.File = ""
	r.Line = 0
	r.OK = false

	return r
}

func putRecord(r *Record) {
	recordPool.Put(r)
}

func (r *Record) Tag(tag string) *Record {
	if r == nil {
		return r
	}
	if !r.ExistTag(tag) {
		r.Tags = append(r.Tags, tag)
	}

	return r
}

func (r *Record) ExistTag(tag string) bool {
	if r == nil {
		return false
	}

	for _, t := range r.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (r *Record) Fields(fields map[string]interface{}) *Record {
	if r == nil {
		return r
	}

	if r.FieldMap == nil {
		r.FieldMap = make(map[string]interface{})
	}

	for k, v := range fields {
		r.FieldMap[k] = v
	}
	return r
}

func (r *Record) Extra(extra interface{}) *Record {
	if r == nil {
		return r
	}

	r.ExtraData = extra
	return r
}

func (r *Record) Msg(msg string, args ...interface{}) {
	if r == nil {
		return
	}

	r.Time = time.Now()
	r.Message = msg
	r.Args = args

	// we can not putRecord() after handleRecord() because Record's logger or it's parent logger can be cachedLogger
	r.Logger.handleRecord(r)
}

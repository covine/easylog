package easylog

import (
	"container/list"
)

// IFilter mockery --name=IFilter --inpackage --case=underscore --filename=filter_mock.go --structname MockFilter
type IFilter interface {
	Filter(*Event) bool
}

type IFilters interface {
	AddFilter(IFilter)
	RemoveFilter(IFilter)
	Filter(*Event) bool
}

// not thread safe, set filters during init
type Filters struct {
	filters *list.List
}

func newFilters() *Filters {
	return &Filters{
		filters: list.New(),
	}
}

func (f *Filters) AddFilter(fi IFilter) {
	if fi == nil {
		return
	}

	if f.filters == nil {
		f.filters = list.New()
	}

	find := false
	for ele := f.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter == fi {
			find = true
			break
		}
	}

	if find {
		return
	} else {
		f.filters.PushBack(fi)
	}
}

func (f *Filters) RemoveFilter(fi IFilter) {
	if fi == nil {
		return
	}

	if f.filters == nil {
		f.filters = list.New()
	}

	var next *list.Element
	for ele := f.filters.Front(); ele != nil; ele = next {
		next = ele.Next()
		filter, ok := ele.Value.(IFilter)
		if ok && filter == fi {
			f.filters.Remove(ele)
		}
	}
}

func (f *Filters) Filter(record *Event) bool {
	if f.filters == nil {
		return true
	}

	for ele := f.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter != nil {
			if filter.Filter(record) == false {
				return false
			}
		}
	}
	return true
}

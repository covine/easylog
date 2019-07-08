package easylog

import (
	"strconv"
	"sync"
	"testing"
)

/*
type LoggerGenerator func() ILogger

func defaultLoggerGenerator() ILogger {
	return &Logger{}
}

func (m *manager) SetLoggerGenerator(g LoggerGenerator) {
	m.loggerGenerator = g
}
*/

func TestLog(t *testing.T) {

	t.Run("get-logger", func(t *testing.T) {
		root := GetRootLogger()
		if root.name != "root" {
			t.Errorf("error root name")
		}

		empty := GetLogger("")
		if empty.parent != root {
			t.Errorf("empty string parent not root")
		}

		empty_empty := GetLogger(".")
		if empty_empty.parent != empty {
			t.Errorf("empty_empty string parent not empty")
		}

		empty_a := GetLogger(".a")
		if empty_a.parent != empty {
			t.Errorf("empty_a string parent not empty")
		}

		empty_empty_a := GetLogger("..a")
		if empty_empty_a.parent != empty_empty {
			t.Errorf("empty_empty_a string parent not empty_empty")
		}

		empty_empty_a_empty_empty := GetLogger("..a..")
		if empty_empty_a_empty_empty.parent != empty_empty_a {
			t.Errorf("empty_empty_a_empty_empty string parent not empty_empty_a")
		}
		if empty_empty_a_empty_empty.parent != GetLogger("..a.") {
			t.Errorf("empty_empty_a_empty_empty string parent not empty_empty_a_empty")
		}

		a_b_c_d_e := GetLogger("a.b.c.d.e")
		if a_b_c_d_e.parent != GetLogger("a") {
			t.Errorf("error parent")
		}

		a_b := GetLogger("a.b")
		if a_b_c_d_e.parent != a_b {
			t.Errorf("error parent")
		}

		a_b_c_d := GetLogger("a.b.c.d")
		if a_b_c_d_e.parent == GetLogger("a.b") {
			t.Errorf("error parent")
		}

		if a_b_c_d_e.parent != a_b_c_d {
			t.Errorf("error parent")
		}

		if a_b_c_d.parent != a_b {
			t.Errorf("error parent")
		}

		if a_b_c_d_e.parent != GetLogger("a.b.c.d") {
			t.Errorf("error parent")
		}

		a_b_c_d_e_d_c := GetLogger("a.b.c.d.e.d.c")
		if a_b_c_d_e_d_c.parent != a_b_c_d_e {
			t.Errorf("error parent")
		}

		b_b_c_d_e_d_c := GetLogger("b.b.c.d.e.d.c")
		if b_b_c_d_e_d_c.parent != root {
			t.Errorf("error parent")
		}
	})

	t.Run("concurrent-get-logger", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10000; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()
				l := GetLogger(strconv.Itoa(j))
				if l.parent != GetRootLogger() {
					t.Errorf("logger parent not root")
				}
			}(i)
		}
		w.Wait()
	})
}

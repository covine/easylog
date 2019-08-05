package easylog

import (
	"strconv"
	"sync"
	"testing"
)

func TestGetLogger(t *testing.T) {

	t.Run("get-logger", func(t *testing.T) {
		root := GetRootLogger()
		if root.name != "root" {
			t.Errorf("error root name")
		}

		empty := GetLogger("")
		if empty.parent != root {
			t.Errorf("empty string parent not root")
		}

		emptyEmpty := GetLogger(".")
		if emptyEmpty.parent != empty {
			t.Errorf("emptyEmpty string parent not empty")
		}

		emptyA := GetLogger(".a")
		if emptyA.parent != empty {
			t.Errorf("emptyA string parent not empty")
		}

		emptyEmptyA := GetLogger("..a")
		if emptyEmptyA.parent != emptyEmpty {
			t.Errorf("emptyEmptyA string parent not emptyEmpty")
		}

		emptyEmptyAEmptyEmpty := GetLogger("..a..")
		if emptyEmptyAEmptyEmpty.parent != emptyEmptyA {
			t.Errorf("emptyEmptyAEmptyEmpty string parent not emptyEmptyA")
		}
		if emptyEmptyAEmptyEmpty.parent != GetLogger("..a.") {
			t.Errorf("emptyEmptyAEmptyEmpty string parent not empty_empty_a_empty")
		}

		a5 := GetLogger("a.b.c.d.e")
		if a5.parent != GetLogger("a") {
			t.Errorf("error parent")
		}

		AB := GetLogger("a.b")
		if a5.parent != AB {
			t.Errorf("error parent")
		}

		a4 := GetLogger("a.b.c.d")
		if a5.parent == GetLogger("a.b") {
			t.Errorf("error parent")
		}

		if a5.parent != a4 {
			t.Errorf("error parent")
		}

		if a4.parent != AB {
			t.Errorf("error parent")
		}

		if a5.parent != GetLogger("a.b.c.d") {
			t.Errorf("error parent")
		}

		a7 := GetLogger("a.b.c.d.e.d.c")
		if a7.parent != a5 {
			t.Errorf("error parent")
		}

		b7 := GetLogger("b.b.c.d.e.d.c")
		if b7.parent != root {
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

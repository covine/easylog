package easylog

import "testing"

func TestLevel(t *testing.T) {
	t.Run("check-level", func(t *testing.T) {
		if NOTSET != 0 {
			t.Errorf("NOTSET is not 0")
		}
		println("NOTSET", NOTSET)

		if DEBUG != 1 {
			t.Errorf("DEBUG is not 1")
		}
		println("DEBUG", DEBUG)

		if INFO != 2 {
			t.Errorf("INFO is not 2")
		}
		println("INFO", INFO)

		if WARNING != 3 {
			t.Errorf("WARNING is not 3")
		}
		println("WARNING", WARNING)

		if WARN != 3 {
			t.Errorf("WARN is not 3")
		}
		println("WARN", WARN)

		if ERROR != 4 {
			t.Errorf("ERROR is not 4")
		}
		println("ERROR", ERROR)

		if FATAL != 5 {
			t.Errorf("FATAL is not 5")
		}
		println("FATAL", FATAL)
	})
}

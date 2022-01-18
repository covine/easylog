package easylog

import "os"

var exit = os.Exit

func modifyExit(f *fakeExit) {
	exit = f.exit
}

func recoverExit() {
	exit = os.Exit
}

type fakeExit struct {
	c      int
	exited bool
}

func (f *fakeExit) exit(code int) {
	f.c = code
	f.exited = true
}

func (f *fakeExit) code() int {
	return f.c
}

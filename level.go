package easylog

type Level uint32

const (
	NOTSET  Level     = iota // 0
	DEBUG                    // 1
	INFO                     // 2
	WARNING                  // 3
	ERROR                    // 4
	FATAL                    // 5
	WARN    = WARNING        // 3
)

func IsLevel(level Level) bool {
	return !(level != NOTSET && level != DEBUG && level != INFO && level != ERROR && level != FATAL && level != WARN)
}

package easylog

type Record struct {
	Level Level
	Msg   string
	Args  []interface{}
}

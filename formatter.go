package easylog

type IFormatter interface {
	Format(record Record) string
}

package handler

import (
	"runtime"
)

var (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Faint      = "\033[2m"
	Underlined = "\033[4m"
	Blink      = "\033[5m"

	Black  = "\033[30m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"

	White = "\033[97m"

	GrayBlack   = "\033[47;30m"
	CyanRed     = "\033[46;31m"
	PurpleGreen = "\033[45;32m"
	BlueYellow  = "\033[44;33m"
	YellowBlue  = "\033[43;34m"
	GreenPurple = "\033[42;35m"
	RedCyan     = "\033[41;36m"
	BlackGray   = "\033[40;37m"

	BlackWhite = "\033[40;97m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Bold = ""
		Faint = ""
		Underlined = ""
		Blink = ""

		Black = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""

		White = ""

		GrayBlack = ""
		CyanRed = ""
		PurpleGreen = ""
		BlueYellow = ""
		YellowBlue = ""
		GreenPurple = ""
		RedCyan = ""
		BlackGray = ""

		BlackWhite = ""
	}
}

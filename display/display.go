package display

import (
	"fmt"
	"mud/interfaces"
)

const (
	Reset         = "\033[0m"
	Black         = "\033[30m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Magenta       = "\033[35m"
	Cyan          = "\033[36m"
	White         = "\033[37m"
	BrightBlack   = "\033[30;1m"
	BrightRed     = "\033[31;1m"
	BrightGreen   = "\033[32;1m"
	BrightYellow  = "\033[33;1m"
	BrightBlue    = "\033[34;1m"
	BrightMagenta = "\033[35;1m"
	BrightCyan    = "\033[36;1m"
	BrightWhite   = "\033[37;1m"
)

func Colorize(text string, color string) string {
	return color + text + Reset
}

func PrintWithColor(player interfaces.PlayerInterface, text string, colorUse string) {
	colorProfile := player.GetColorProfile()
	color := colorProfile.GetColor(colorUse)
	playerConn := player.GetConn()
	fmt.Fprintf(playerConn, "%s", Colorize(text, color))
}

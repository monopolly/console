package console

import "github.com/fatih/color"

var (
	Yellow = color.New(color.FgHiYellow).SprintFunc()
	Red    = color.New(color.FgHiRed).SprintFunc()
	Cyan   = color.New(color.FgHiCyan).SprintFunc()
	Mag    = color.New(color.FgHiMagenta).SprintFunc()
	Green  = color.New(color.FgHiGreen).SprintFunc()
	White  = color.New(color.FgWhite).SprintFunc()
	Black  = color.New(color.FgHiBlack).SprintFunc()
	Blue   = color.New(color.FgHiBlue).SprintFunc()
	Gray   = color.New(color.FgHiWhite).SprintFunc()
)

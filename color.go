package treelang

import "github.com/fatih/color"

// // // // // // // // // //

var (
	Red     = color.New(color.FgHiRed).SprintFunc()
	Green   = color.New(color.FgHiGreen).SprintFunc()
	Yellow  = color.New(color.FgHiYellow).SprintFunc()  // h3
	Blue    = color.New(color.FgHiBlue).SprintFunc()    // h2
	Magenta = color.New(color.FgHiMagenta).SprintFunc() // h1
	Cyan    = color.New(color.FgHiCyan).SprintFunc()    // i
)

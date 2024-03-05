
package utils

import (
	"fmt"
	"grep/worker"
)

var (
	// ANSI escape codes for text colors
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)


func PrintColoredResult(r worker.Result) {

	fmt.Printf("%s%s > found %s'%s' at line %s%d%s\n",
	colorGreen, r.Path, colorYellow, r.Line, colorRed, r.LineNum, colorReset)
}
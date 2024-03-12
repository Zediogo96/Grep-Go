package utils

import (
	"fmt"
	"grep/worker"
	"strings"
)

var (
	// ANSI escape codes for text colors
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

type FormattedResult struct {
	Line         string
	diffStartPos int
	diffEndPos   int
}

func PrintColoredResult(r worker.Result, pattern string) {

	if r.Line == "" {
		return
	}

	var formattedResult FormattedResult

	formattedResult.Line = r.Line
	formattedResult.diffStartPos = strings.Index(strings.ToLower(r.Line), strings.ToLower(pattern))
	formattedResult.diffEndPos = formattedResult.diffStartPos + len(pattern)

	// If the pattern is not found, print the line as is, and return
	if formattedResult.diffStartPos == -1 {
		return
	}

	beforePattern := r.Line[:formattedResult.diffStartPos]
	afterPattern := r.Line[formattedResult.diffEndPos:]
	patternColor := colorRed + pattern + colorReset
	pathColor := colorYellow + r.Path + colorReset
	lineNumColor := colorGreen + fmt.Sprint(r.LineNum) + colorReset

	// Print the string, using FormattedResult and
	fmt.Printf("%s > line %s > %s%s%s \n", pathColor, lineNumColor, beforePattern, patternColor, afterPattern)

}

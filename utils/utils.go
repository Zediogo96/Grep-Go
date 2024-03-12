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
	diffStartPos int
	diffEndPos   int
}

func PrintResultsColored(r worker.Result, pattern string) {

	if r.Line == "" {
		return
	}

	line := r.Line

	var formattedResult FormattedResult
	startPos := 0
	matchedPositions := []FormattedResult{} // Keep track of matched positions

	absoluteIndex := 0

	for {
		startPos = strings.Index(line, pattern)

		if startPos == -1 {
			absoluteIndex = 0
			break
		}

		endPos := startPos + len(pattern)

		formattedResult.diffStartPos = startPos + absoluteIndex
		formattedResult.diffEndPos = endPos + absoluteIndex
		matchedPositions = append(matchedPositions, formattedResult)

		line = line[endPos:]
		// Remove the matched pattern from the line
		absoluteIndex += endPos
	}

	coloredPath := colorYellow + r.Path + colorReset
	coloredLineNum := colorGreen + fmt.Sprintf("%d", r.LineNum) + colorReset

	fmt.Printf("%s > line n. %s > ", coloredPath, coloredLineNum)
	for c := range r.Line {
		for _, pos := range matchedPositions {

			if c == pos.diffStartPos {
				fmt.Print(colorRed)
			} else if c == pos.diffEndPos {
				fmt.Print(colorReset)
			}
		}
		fmt.Print(string(r.Line[c]))
	}

	fmt.Println()
}

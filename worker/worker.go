package worker

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Result struct {
    Line    string
    LineNum int
    Path    string
}

type Results struct {
    Inner []Result
}
func NewResult(line string, lineNum int, path string) Result {
    return Result{line, lineNum, path}
}

func ProcessFile(path string, stringToFind string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file ", path)
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var results Results

	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, stringToFind) {
			results.Inner = append(results.Inner, NewResult(line, lineNum, path))
		}
		lineNum++
	}

	return &results
}
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

func ProcessFile(path string, stringToFind string, caseSensitive bool) *Results {
	file, err := os.Open(path)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }
    results := Results{make([]Result, 0)}
    scanner := bufio.NewScanner(file)
    lineNum := 1
    for scanner.Scan() {
        line := scanner.Text()
        if caseSensitive {
            if strings.Contains(line, stringToFind) {
                results.Inner = append(results.Inner, NewResult(line, lineNum, path))
            }
        } else {
            if strings.Contains(strings.ToLower(line), strings.ToLower(stringToFind)) {
                results.Inner = append(results.Inner, NewResult(line, lineNum, path))
            }
        }
        lineNum++
    }
    if len(results.Inner) == 0 {
        return nil
    } else {
        return &results
    }
}
package worker

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func ProcessFile(path string, stringToFind string, caseSensitive bool, useRegex bool) *Results {
	file, err := os.Open(path)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }
    results := Results{make([]Result, 0)}
    scanner := bufio.NewScanner(file)
    lineNum := 1
    for scanner.Scan() {
        if useRegex {

            re, err := regexp.Compile(stringToFind)

            if err != nil {
                fmt.Println("Error:", err)
                return nil
            }
            if re.MatchString(scanner.Text()) {
                results.Inner = append(results.Inner, NewResult(scanner.Text(), lineNum, path))
            }
        } else {
            if caseSensitive {
                if strings.Contains(scanner.Text(), stringToFind) {
                    results.Inner = append(results.Inner, NewResult(scanner.Text(), lineNum, path))
                }
            } else {
                if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(stringToFind)) {
                    results.Inner = append(results.Inner, NewResult(scanner.Text(), lineNum, path))
                }
            }
        }

    }
    if len(results.Inner) == 0 {
        return nil
    } else {
        return &results
    }
}

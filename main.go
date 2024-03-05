package main

import (
	"fmt"
	"os"
	"grep/worklist"
)


func extractAllFiles(wl *worklist.Worklist, path string) {

	// Get all the files in the current directory
	entries, err := os.ReadDir(path)

	// An error occurred with os.ReadDir
	if err != nil {
		fmt.Println("Error reading directory ", path)
		return
	}


	// Iterate through all the files in the current directory
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("Adding directory ", entry.Name())
			nextPath := path + "/" + entry.Name()
			extractAllFiles(wl, nextPath)

		// Meaning that the entry is a file
		} else {
			fmt.Println("Adding file ", entry.Name())
			wl.Add(worklist.NewJob(path + "/" + entry.Name()))
		}
	}
}

func main() {
	wl := worklist.New(100)
	numWorkers := 10

	path := os.Args[2]

	extractAllFiles(&wl, path)
	wl.Finalize(numWorkers)

}

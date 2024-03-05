package main

import (
	"fmt"
	"os"
)

func main() {
	// wl := worklist.New(100)

	// numWorkers := 10

	path := os.Args[2]

	// Get all the files in the current directory
	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Error reading directory ", path)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

}
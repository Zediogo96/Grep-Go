package main

import (
	"flag"
	"fmt"
	"grep/utils"
	"grep/worker"
	"grep/worklist"
	"os"
	"sync"
)

func extractAllFiles(wl *worklist.Worklist, path string) {

	// * Get all the files in the current directory
	entries, err := os.ReadDir(path)

	// * An error occurred with os.ReadDir
	if err != nil {
		fmt.Println("Error reading directory ", path)
		return
	}

	// * Iterate through all the files in the current directory
	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := path + entry.Name()
			extractAllFiles(wl, nextPath)

			// * Meaning that the entry is a file
		} else {
			wl.Add(worklist.NewJob(path + "/" + entry.Name()))
		}
	}
}

func main() {
	var workersWG sync.WaitGroup

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	workersWG.Add(1)

	caseSensitive := flag.Bool("i", false, "Perform a case-sensitive search")
	numWorkers := flag.Int("w", 10, "Number of workers to use")
	useRegex := flag.Bool("e", false, "Use a regular expression search")
	flag.Parse()

	// * Print all the active flags
	fmt.Println("Active flags: ")
	flag.PrintDefaults()

	NON_FLAG_ARGS := flag.Args()

	if len(NON_FLAG_ARGS) < 2 {
		fmt.Println("Usage: grep [flags] pattern path")
		os.Exit(1)
	}

	var searchTerm string

	path := NON_FLAG_ARGS[1]

	searchTerm = NON_FLAG_ARGS[0]

	go func() {
		defer workersWG.Done()
		extractAllFiles(&wl, path)
		wl.Finalize(*numWorkers)
	}()

	for i := 0; i < *numWorkers; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			for {
				entry := wl.Next()
				if entry.Path != "" {
					workerResult := worker.ProcessFile(entry.Path, searchTerm, *caseSensitive, *useRegex)
					if workerResult != nil {
						for _, r := range workerResult.Inner {
							results <- r
						}
					}
				} else {
					return // Terminate the worker goroutine when no more files to process
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})
	go func() {
		workersWG.Wait()
		close(blockWorkersWg)
	}()

	var displayWg sync.WaitGroup
	displayWg.Add(1)

	fmt.Println("\nResults: ")

	go func() {
		for {
			select {
			case r := <-results:
				utils.PrintColoredResult(r, searchTerm)
			case <-blockWorkersWg:
				// Make sure channel is empty before aborting display goroutine
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			}
		}
	}()
	displayWg.Wait()
}

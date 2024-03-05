package main

import (
	"fmt"
	"os"
	"grep/worklist"
	"grep/worker"
	"grep/utils"
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
			// fmt.Println("Adding directory ", entry.Name())
			nextPath := path + "/" + entry.Name()
			extractAllFiles(wl, nextPath)

		// * Meaning that the entry is a file
		} else {
			// fmt.Println("Adding file ", entry.Name())
			wl.Add(worklist.NewJob(path + "/" + entry.Name()))
		}
	}
}

func main() {
	var workersWG sync.WaitGroup

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := 10

	path := os.Args[2]

	workersWG.Add(1)

	go func() {
		defer workersWG.Done()
		extractAllFiles(&wl, path)
		wl.Finalize(numWorkers)
	}()

	for i := 0; i < numWorkers; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			for {
				entry := wl.Next()
				if entry.Path != "" {
					workerResult := worker.ProcessFile(entry.Path, os.Args[1])
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
     // Close channel
     close(blockWorkersWg)
    }()

	var displayWg sync.WaitGroup
    displayWg.Add(1)

	go func() {
        for {
            select {
            case r := <-results:
				utils.PrintColoredResult(r)
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
package main

import (
	"fmt"
	"os"
	"grep/worklist"
	"grep/worker"
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
			fmt.Println("Adding directory ", entry.Name())
			nextPath := path + "/" + entry.Name()
			extractAllFiles(wl, nextPath)

		// * Meaning that the entry is a file
		} else {
			fmt.Println("Adding file ", entry.Name())
			wl.Add(worklist.NewJob(path + "/" + entry.Name()))
		}
	}
}

func main() {


	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := 10

	path := os.Args[2]

	go func() {
		wl.Finalize(numWorkers)
		extractAllFiles(&wl, path)
	}()

	for i := 0; i < numWorkers; i++ {
		go func() {
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
					// ! If the path is empty, then the worker should terminate as there
					// ! are no more files to process
					return
				}

			}
		}()
	}

	// ? Wait for all the workers to finish
	go func() {
		for {
         select {
         case r := <-results:
             fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNum, r.Line)
         }
        }
	}()

}

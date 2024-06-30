package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	numUUIDs   = 10_000_000
	numWorkers = 1_000 // Adjust this number based on your CPU cores and workload
)

func generateUUIDs(id int, numUUIDs int, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()

	for i := 0; i < numUUIDs; i++ {
		_, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("Worker %d: Failed to generate UUID: %v", id, err)
		}
		// Optionally, you can uncomment the following lines to see the progress
		if i%1000000 == 0 {
			fmt.Printf("Worker %d: Generated %d UUIDs, %d CPUs, %d GoRoutines \n", id, i, runtime.NumCPU(), runtime.NumGoroutine())
		}
	}

	ch <- true
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	ch := make(chan bool, numWorkers)

	numUUIDsPerWorker := numUUIDs / numWorkers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go generateUUIDs(i, numUUIDsPerWorker, &wg, ch)
	}

	wg.Wait()
	close(ch)

	elapsed := time.Since(start)
	fmt.Printf("Generated %d UUIDs in %s\n", numUUIDs, elapsed)
}

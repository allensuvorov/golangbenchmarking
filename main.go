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
	numUUIDs = 100_000
	// numWorkers = 1000 // Adjust this number based on your CPU cores and workload
)

func generateUUIDs(id int, numUUIDs int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < numUUIDs; i++ {
		_, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("Worker %d: Failed to generate UUID: %v", id, err)
		}
		// Optionally, you can uncomment the following lines to see the progress
		// if i%1000000 == 0 {
		// 	fmt.Printf("Worker %d: Generated %d UUIDs, %d CPUs, %d GoRoutines \n", id, i, runtime.NumCPU(), runtime.NumGoroutine())
		// }
	}
}

func main() {
	var bestTime int64 = 1_000_000_000
	bestNumWorkers := -1
	bestNumCPU := -1
	for CPUs := 4; CPUs <= 4; CPUs++ {
		for numWorkers := 1; numWorkers <= 100_001; numWorkers += 10_000 {
			prevCPUs := runtime.GOMAXPROCS(CPUs)
			start := time.Now()

			var wg sync.WaitGroup

			numUUIDsPerWorker := numUUIDs / numWorkers

			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				go generateUUIDs(i, numUUIDsPerWorker, &wg)
			}

			wg.Wait()

			elapsed := time.Since(start).Milliseconds()

			if elapsed < bestTime {
				bestTime = elapsed
				bestNumWorkers = numWorkers
				bestNumCPU = CPUs
			}

			fmt.Printf("Generated %d UUIDs in %v milli sec, %d prevCPUs, %d CPUs, %v workers\n", numUUIDs, elapsed, prevCPUs, CPUs, numWorkers)

		}
	}
	fmt.Printf("To generate %d UUIDs, best time is %v milli seconds, with %v usable CPUs and %v workers", numUUIDs, bestTime, bestNumCPU, bestNumWorkers)
}

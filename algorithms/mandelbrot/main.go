package main

import (
	"fmt"
	"math/cmplx"
	"sync"
	"time"
)

func mandelbrot(c complex128, maxIter int) int {
	z := c
	for n := 0; n < maxIter; n++ {
		if cmplx.Abs(z) > 2 {
			return n
		}
		z = z*z + c
	}
	return maxIter
}

func sequentialMandelbrot(width, height, maxIter int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		row := make([]int, width)
		for x := 0; x < width; x++ {
			re := (float64(x)/float64(width))*3.5 - 2.5
			im := (float64(y)/float64(height))*2.0 - 1.0
			row[x] = mandelbrot(complex(re, im), maxIter)
		}
		pixels[y] = row
	}
	return pixels
}

func parallelMandelbrot(width, height, maxIter int) [][]int {
	pixels := make([][]int, height)
	var wg sync.WaitGroup

	for y := 0; y < height; y++ {
		row := make([]int, width)
		wg.Add(1)
		go func(y int, row []int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				re := (float64(x)/float64(width))*3.5 - 2.5
				im := (float64(y)/float64(height))*2.0 - 1.0
				row[x] = mandelbrot(complex(re, im), maxIter)
			}
			pixels[y] = row
		}(y, row)
	}

	wg.Wait()
	return pixels
}

func main() {
	width := 3840  // Increased width 20x (1920 * 2)
	height := 2160 // Increased height 20x (1080 * 2)
	maxIter := 1000

	// Sequential computation
	startSeq := time.Now()
	pixelsSeq := sequentialMandelbrot(width, height, maxIter)
	durationSeq := time.Since(startSeq)

	// Parallel computation
	startPar := time.Now()
	pixelsPar := parallelMandelbrot(width, height, maxIter)
	durationPar := time.Since(startPar)

	// Print the results
	fmt.Printf("Time taken for sequential Mandelbrot set computation: %v\n", durationSeq)
	fmt.Printf("Number of pixels computed sequentially: %d\n", len(pixelsSeq)*len(pixelsSeq[0]))

	fmt.Printf("Time taken for parallel Mandelbrot set computation: %v\n", durationPar)
	fmt.Printf("Number of pixels computed in parallel: %d\n", len(pixelsPar)*len(pixelsPar[0]))
}

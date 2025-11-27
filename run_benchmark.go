package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abola/wa-tor/simulation"
)

type BenchResult struct {
	Threads        int
	TimeMs         float64
	ChrononsPerSec float64
}

func RunBenchmark(threads int) BenchResult {
	gridSize := 30
	numFish := 200
	numSharks := 40
	duration := 50

	rand.Seed(time.Now().UnixNano())
	sim := simulation.NewSimulation(gridSize, numFish, numSharks, 4, 8, 10)

	start := time.Now()
	for i := 0; i < duration; i++ {
		if threads == 1 {
			sim.Step()
		} else {
			sim.StepParallel(threads)
		}
	}

	elapsed := time.Since(start)
	timeMs := elapsed.Seconds() * 1000
	chrononsPerSec := float64(duration) / elapsed.Seconds()

	return BenchResult{Threads: threads, TimeMs: timeMs, ChrononsPerSec: chrononsPerSec}
}

func BenchmarkMain() {
	fmt.Println("=== Wa-Tor Concurrent Multithreading Benchmark ===\n")

	results := []BenchResult{}
	for _, threads := range []int{1, 2, 4, 8} {
		fmt.Printf("Testing %d thread(s)...\n", threads)
		result := RunBenchmark(threads)
		results = append(results, result)
		fmt.Printf("  Time: %.2f ms | Chronons/sec: %.2f\n\n", result.TimeMs, result.ChrononsPerSec)
	}

	fmt.Println("=== Speedup Analysis ===")
	baseline := results[0].ChrononsPerSec
	for _, r := range results {
		fmt.Printf("Threads: %d | Speedup: %.2fx\n", r.Threads, r.ChrononsPerSec/baseline)
	}
}

func main() {
	BenchmarkMain()
}

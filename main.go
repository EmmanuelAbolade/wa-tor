package main

import (
	"flag"
	"fmt"
	"github.com/abola/wa-tor/simulation"
	"math/rand"
	"time"
)

func main() {
	// Command-line flags
	numShark := flag.Int("sharks", 100, "Starting population of sharks")
	numFish := flag.Int("fish", 500, "Starting population of fish")
	fishBreed := flag.Int("fishbreed", 4, "Number of chronons before fish reproduce")
	sharkBreed := flag.Int("sharkbreed", 8, "Number of chronons before shark reproduces")
	starve := flag.Int("starve", 10, "Chronons before shark starves")
	gridSize := flag.Int("grid", 50, "Grid dimensions (size x size)")
	numThreads := flag.Int("threads", 1, "Number of threads (not used yet)")
	duration := flag.Int("duration", 500, "Number of chronons to simulate")

	flag.Parse()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Wa-Tor Simulation ===")
	fmt.Printf("Grid Size: %d x %d\n", *gridSize, *gridSize)
	fmt.Printf("Initial Fish: %d, Sharks: %d\n", *numFish, *numShark)
	fmt.Printf("Fish Breed Age: %d, Shark Breed Age: %d\n", *fishBreed, *sharkBreed)
	fmt.Printf("Shark Starve Time: %d\n", *starve)
	fmt.Printf("Threads: %d\n", *numThreads)
	fmt.Printf("Duration: %d chronons\n\n", *duration)

	// Create simulation
	sim := simulation.NewSimulation(*gridSize, *numFish, *numShark, *fishBreed, *sharkBreed, *starve)

	// Run simulation
	startTime := time.Now()

	for i := 0; i < *duration; i++ {
		sim.Step()

		// Print stats every 50 chronons
		if i%50 == 0 {
			sim.PrintStats()
		}
	}

	elapsed := time.Since(startTime)

	// Final stats
	fmt.Println("\n=== Final Results ===")
	sim.PrintStats()
	fmt.Printf("Time elapsed: %v\n", elapsed)
	fmt.Printf("Chronons per second: %.2f\n", float64(*duration)/elapsed.Seconds())
}

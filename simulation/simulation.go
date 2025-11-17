package simulation

import (
	"fmt"
	"math/rand"
	"sync"
)

// Simulation represents the entire Wa-Tor world
type Simulation struct {
	Grid            *Grid
	Fish            map[int]*Fish
	Sharks          map[int]*Shark
	Chronon         int
	FishBreedAge    int
	SharkBreedAge   int
	SharkStarveTime int
	EnergyPerFish   int
	NextFishID      int
	NextSharkID     int
	mu              sync.RWMutex
}

// NewSimulation initializes a new simulation
func NewSimulation(gridSize, numFish, numSharks, fishBreedAge, sharkBreedAge, sharkStarveTime int) *Simulation {
	sim := &Simulation{
		Grid:            NewGrid(gridSize),
		Fish:            make(map[int]*Fish),
		Sharks:          make(map[int]*Shark),
		Chronon:         0,
		FishBreedAge:    fishBreedAge,
		SharkBreedAge:   sharkBreedAge,
		SharkStarveTime: sharkStarveTime,
		EnergyPerFish:   3,
		NextFishID:      1,
		NextSharkID:     1,
	}

	// Place initial fish randomly
	for i := 0; i < numFish; i++ {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)

		// Find an empty spot
		for sim.Grid.GetCell(x, y).Type != EMPTY {
			x = rand.Intn(gridSize)
			y = rand.Intn(gridSize)
		}

		fish := NewFish(sim.NextFishID, x, y, fishBreedAge)
		sim.NextFishID++
		sim.Fish[fish.ID] = fish
		sim.Grid.SetCell(x, y, Cell{Type: FISH, ID: fish.ID})
	}

	// Place initial sharks randomly
	for i := 0; i < numSharks; i++ {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)

		// Find an empty spot
		for sim.Grid.GetCell(x, y).Type != EMPTY {
			x = rand.Intn(gridSize)
			y = rand.Intn(gridSize)
		}

		shark := NewShark(sim.NextSharkID, x, y, sharkStarveTime)
		sim.NextSharkID++
		sim.Sharks[shark.ID] = shark
		sim.Grid.SetCell(x, y, Cell{Type: SHARK, ID: shark.ID})
	}

	return sim
}

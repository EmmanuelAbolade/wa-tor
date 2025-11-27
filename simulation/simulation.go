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
	mapMutex        sync.Mutex
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

// Step advances the simulation by one chronon using sequential processing.
// It processes all fish and sharks in turn, updating their positions,
// energy levels, and reproducing as appropriate.
func (s *Simulation) Step() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stepFish()
	s.stepSharks()

	s.Chronon++
}

// stepFish handles all fish movement and reproduction
func (s *Simulation) stepFish() {
	fishToProcess := make([]*Fish, 0, len(s.Fish))

	// Get current list of fish
	for _, fish := range s.Fish {
		fishToProcess = append(fishToProcess, fish)
	}

	for _, fish := range fishToProcess {
		if _, exists := s.Fish[fish.ID]; !exists {
			continue // Fish was eaten or removed
		}

		// Get adjacent empty cells
		empty, _, _ := s.Grid.GetAdjacentCells(fish.X, fish.Y)

		// Move fish if there's space
		if len(empty) > 0 {
			// Choose random empty adjacent cell
			newPos := empty[rand.Intn(len(empty))]

			// Remove fish from old position
			s.Grid.SetCell(fish.X, fish.Y, Cell{Type: EMPTY, ID: 0})

			// Move fish to new position
			fish.X = newPos.x
			fish.Y = newPos.y
			s.Grid.SetCell(fish.X, fish.Y, Cell{Type: FISH, ID: fish.ID})
		}

		// Increase age
		fish.IncreaseAge()

		// Check if fish can reproduce
		if fish.CanReproduce(s.FishBreedAge) {
			// Create offspring at old position
			newFish := NewFish(s.NextFishID, fish.X, fish.Y, s.FishBreedAge)
			s.NextFishID++
			s.Fish[newFish.ID] = newFish
			s.Grid.SetCell(newFish.X, newFish.Y, Cell{Type: FISH, ID: newFish.ID})

			// Reset parent's age
			fish.Age = 0
		}
	}
}

// stepSharks handles all shark movement, hunting, starving, and reproduction
func (s *Simulation) stepSharks() {
	sharksToProcess := make([]*Shark, 0, len(s.Sharks))

	// Get current list of sharks
	for _, shark := range s.Sharks {
		sharksToProcess = append(sharksToProcess, shark)
	}

	for _, shark := range sharksToProcess {
		if _, exists := s.Sharks[shark.ID]; !exists {
			continue // Shark already processed/dead
		}

		// Get adjacent cells
		empty, fish, _ := s.Grid.GetAdjacentCells(shark.X, shark.Y)

		// Step 1: Sharks prioritize eating fish
		if len(fish) > 0 {
			// Move to a random fish
			preyPos := fish[rand.Intn(len(fish))]
			preyCell := s.Grid.GetCell(preyPos.x, preyPos.y)

			// Eat the fish
			if preyFish, exists := s.Fish[preyCell.ID]; exists {
				delete(s.Fish, preyFish.ID)
			}

			// Move shark to fish position
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			shark.X = preyPos.x
			shark.Y = preyPos.y
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: SHARK, ID: shark.ID})

			// Eat gives energy
			shark.Eat(s.EnergyPerFish)
		} else if len(empty) > 0 {
			// Step 2: No fish nearby, move to empty space (like a fish)
			newPos := empty[rand.Intn(len(empty))]

			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			shark.X = newPos.x
			shark.Y = newPos.y
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: SHARK, ID: shark.ID})
		}

		// Step 3: Shark loses energy each chronon
		shark.Starve()
		shark.IncreaseAge()

		// Step 4: Check if shark starves
		if !shark.IsAlive() {
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			delete(s.Sharks, shark.ID)
			continue
		}

		// Step 5: Check if shark can reproduce
		if shark.CanReproduce(s.SharkBreedAge) {
			// Create offspring at old position
			newShark := NewShark(s.NextSharkID, shark.X, shark.Y, shark.Energy/2)
			s.NextSharkID++
			s.Sharks[newShark.ID] = newShark
			s.Grid.SetCell(newShark.X, newShark.Y, Cell{Type: SHARK, ID: newShark.ID})

			// Parent's energy splits
			shark.Energy = shark.Energy / 2
			shark.Age = 0
		}
	}
}

// GetStats returns current population counts
func (s *Simulation) GetStats() (chronon int, numFish int, numSharks int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Chronon, len(s.Fish), len(s.Sharks)
}

// PrintStats prints current population statistics
func (s *Simulation) PrintStats() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Printf("Chronon: %d | Fish: %d | Sharks: %d\n", s.Chronon, len(s.Fish), len(s.Sharks))
}

// GetGridCopy returns a snapshot of the grid (for rendering)
func (s *Simulation) GetGridCopy() [][]Cell {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make([][]Cell, s.Grid.Size)
	for i := range copy {
		copy[i] = make([]Cell, s.Grid.Size)
		for j := range copy[i] {
			copy[i][j] = s.Grid.Cells[i][j]
		}
	}
	return copy
}

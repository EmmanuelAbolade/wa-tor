package simulation

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Simulation represents the entire Wa-Tor world
type Simulation struct {
	Grid            *Grid
	Fish            sync.Map // Thread-safe map for fish
	Sharks          sync.Map // Thread-safe map for sharks
	Chronon         int
	FishBreedAge    int
	SharkBreedAge   int
	SharkStarveTime int
	EnergyPerFish   int
	NextFishID      int
	NextSharkID     int
	FishCount       atomic.Int64 // Thread-safe counter for fish
	SharkCount      atomic.Int64 // Thread-safe counter for sharks
	mu              sync.RWMutex // Protects Chronon
}

// NewSimulation initializes a new simulation
func NewSimulation(gridSize, numFish, numSharks, fishBreedAge, sharkBreedAge, sharkStarveTime int) *Simulation {
	sim := &Simulation{
		Grid:            NewGrid(gridSize),
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
		sim.Fish.Store(fish.ID, fish) // Use Store for sync.Map
		sim.FishCount.Add(1)          // Increment counter
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
		sim.Sharks.Store(shark.ID, shark) // Use Store for sync.Map
		sim.SharkCount.Add(1)             // Increment counter
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
	fishToProcess := make([]*Fish, 0)

	// Get current list of fish using Range (thread-safe iteration)
	s.Fish.Range(func(key, value interface{}) bool {
		fish := value.(*Fish)
		fishToProcess = append(fishToProcess, fish)
		return true
	})

	for _, fish := range fishToProcess {
		if _, exists := s.Fish.Load(fish.ID); !exists {
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
			s.Fish.Store(newFish.ID, newFish) // Use Store for sync.Map
			s.FishCount.Add(1)
			s.Grid.SetCell(newFish.X, newFish.Y, Cell{Type: FISH, ID: newFish.ID})

			// Reset parent's age
			fish.Age = 0
		}
	}
}

// stepSharks handles all shark movement, hunting, starving, and reproduction
func (s *Simulation) stepSharks() {
	sharksToProcess := make([]*Shark, 0)

	// Get current list of sharks using Range (thread-safe iteration)
	s.Sharks.Range(func(key, value interface{}) bool {
		shark := value.(*Shark)
		sharksToProcess = append(sharksToProcess, shark)
		return true
	})

	for _, shark := range sharksToProcess {
		if _, exists := s.Sharks.Load(shark.ID); !exists {
			continue // Shark already processed/dead
		}

		// Get adjacent cells
		empty, fish, _ := s.Grid.GetAdjacentCells(shark.X, shark.Y)

		// Step 1: Sharks prioritize eating fish
		if len(fish) > 0 {
			// Move to a random fish
			preyPos := fish[rand.Intn(len(fish))]
			preyCell := s.Grid.GetCell(preyPos.x, preyPos.y)

			// Eat the fish (thread-safe delete)
			if _, exists := s.Fish.Load(preyCell.ID); exists {
				s.Fish.Delete(preyCell.ID) // Use Delete for sync.Map
				s.FishCount.Add(-1)
			}

			// Move shark to fish position
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			shark.X = preyPos.x
			shark.Y = preyPos.y
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: SHARK, ID: shark.ID})

			// Eat gives energy
			shark.Eat(s.EnergyPerFish)
		} else if len(empty) > 0 {
			// Step 2: No fish nearby, move to empty space
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
			s.Sharks.Delete(shark.ID) // Use Delete for sync.Map
			s.SharkCount.Add(-1)
			continue
		}

		// Step 5: Check if shark can reproduce
		if shark.CanReproduce(s.SharkBreedAge) {
			// Create offspring at old position
			newShark := NewShark(s.NextSharkID, shark.X, shark.Y, shark.Energy/2)
			s.NextSharkID++
			s.Sharks.Store(newShark.ID, newShark) // Use Store for sync.Map
			s.SharkCount.Add(1)
			s.Grid.SetCell(newShark.X, newShark.Y, Cell{Type: SHARK, ID: newShark.ID})

			// Parent's energy splits
			shark.Energy = shark.Energy / 2
			shark.Age = 0
		}
	}
}

// GetStats returns current population counts (thread-safe)
func (s *Simulation) GetStats() (chronon int, numFish int, numSharks int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Chronon, int(s.FishCount.Load()), int(s.SharkCount.Load())
}

// PrintStats prints current population statistics
func (s *Simulation) PrintStats() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Printf("Chronon: %d | Fish: %d | Sharks: %d\n", s.Chronon, int(s.FishCount.Load()), int(s.SharkCount.Load()))
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

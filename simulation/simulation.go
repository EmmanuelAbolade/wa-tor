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
	Fish            sync.Map // Thread-safe map
	Sharks          sync.Map // Thread-safe map
	Chronon         int
	FishBreedAge    int
	SharkBreedAge   int
	SharkStarveTime int
	EnergyPerFish   int
	NextFishID      int
	NextSharkID     int
	FishCount       atomic.Int64 // Thread-safe counter
	SharkCount      atomic.Int64 // Thread-safe counter
	Mu              sync.RWMutex // Exported mutex
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

		for sim.Grid.GetCell(x, y).Type != EMPTY {
			x = rand.Intn(gridSize)
			y = rand.Intn(gridSize)
		}

		fish := NewFish(sim.NextFishID, x, y, fishBreedAge)
		sim.NextFishID++
		sim.Fish.Store(fish.ID, fish)
		sim.FishCount.Add(1)
		sim.Grid.SetCell(x, y, Cell{Type: FISH, ID: fish.ID})
	}

	// Place initial sharks randomly
	for i := 0; i < numSharks; i++ {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)

		for sim.Grid.GetCell(x, y).Type != EMPTY {
			x = rand.Intn(gridSize)
			y = rand.Intn(gridSize)
		}

		shark := NewShark(sim.NextSharkID, x, y, sharkStarveTime)
		sim.NextSharkID++
		sim.Sharks.Store(shark.ID, shark)
		sim.SharkCount.Add(1)
		sim.Grid.SetCell(x, y, Cell{Type: SHARK, ID: shark.ID})
	}

	return sim
}

// Step advances the simulation by one chronon
func (s *Simulation) Step() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.stepFish()
	s.stepSharks()

	s.Chronon++
}

// stepFish handles all fish movement and reproduction
func (s *Simulation) stepFish() {
	fishToProcess := make([]*Fish, 0)

	s.Fish.Range(func(key, value interface{}) bool {
		fish := value.(*Fish)
		fishToProcess = append(fishToProcess, fish)
		return true
	})

	for _, fish := range fishToProcess {
		if _, exists := s.Fish.Load(fish.ID); !exists {
			continue
		}

		empty, _, _ := s.Grid.GetAdjacentCells(fish.X, fish.Y)

		if len(empty) > 0 {
			newPos := empty[rand.Intn(len(empty))]
			s.Grid.SetCell(fish.X, fish.Y, Cell{Type: EMPTY, ID: 0})
			fish.X = newPos.x
			fish.Y = newPos.y
			s.Grid.SetCell(fish.X, fish.Y, Cell{Type: FISH, ID: fish.ID})
		}

		fish.IncreaseAge()

		if fish.CanReproduce(s.FishBreedAge) {
			newFish := NewFish(s.NextFishID, fish.X, fish.Y, s.FishBreedAge)
			s.NextFishID++
			s.Fish.Store(newFish.ID, newFish)
			s.FishCount.Add(1)
			s.Grid.SetCell(newFish.X, newFish.Y, Cell{Type: FISH, ID: newFish.ID})
			fish.Age = 0
		}
	}
}

// stepSharks handles all shark movement, hunting, starving, and reproduction
func (s *Simulation) stepSharks() {
	sharksToProcess := make([]*Shark, 0)

	s.Sharks.Range(func(key, value interface{}) bool {
		shark := value.(*Shark)
		sharksToProcess = append(sharksToProcess, shark)
		return true
	})

	for _, shark := range sharksToProcess {
		if _, exists := s.Sharks.Load(shark.ID); !exists {
			continue
		}

		empty, fish, _ := s.Grid.GetAdjacentCells(shark.X, shark.Y)

		if len(fish) > 0 {
			preyPos := fish[rand.Intn(len(fish))]
			preyCell := s.Grid.GetCell(preyPos.x, preyPos.y)

			if _, exists := s.Fish.Load(preyCell.ID); exists {
				s.Fish.Delete(preyCell.ID)
				s.FishCount.Add(-1)
			}

			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			shark.X = preyPos.x
			shark.Y = preyPos.y
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: SHARK, ID: shark.ID})
			shark.Eat(s.EnergyPerFish)
		} else if len(empty) > 0 {
			newPos := empty[rand.Intn(len(empty))]
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			shark.X = newPos.x
			shark.Y = newPos.y
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: SHARK, ID: shark.ID})
		}

		shark.Starve()
		shark.IncreaseAge()

		if !shark.IsAlive() {
			s.Grid.SetCell(shark.X, shark.Y, Cell{Type: EMPTY, ID: 0})
			s.Sharks.Delete(shark.ID)
			s.SharkCount.Add(-1)
			continue
		}

		if shark.CanReproduce(s.SharkBreedAge) {
			newShark := NewShark(s.NextSharkID, shark.X, shark.Y, shark.Energy/2)
			s.NextSharkID++
			s.Sharks.Store(newShark.ID, newShark)
			s.SharkCount.Add(1)
			s.Grid.SetCell(newShark.X, newShark.Y, Cell{Type: SHARK, ID: newShark.ID})
			shark.Energy = shark.Energy / 2
			shark.Age = 0
		}
	}
}

// GetStats returns current population counts
func (s *Simulation) GetStats() (chronon int, numFish int, numSharks int) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Chronon, int(s.FishCount.Load()), int(s.SharkCount.Load())
}

// PrintStats prints current population statistics
func (s *Simulation) PrintStats() {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	fmt.Printf("Chronon: %d | Fish: %d | Sharks: %d\n", s.Chronon, int(s.FishCount.Load()), int(s.SharkCount.Load()))
}

// GetGridCopy returns a snapshot of the grid (for rendering)
func (s *Simulation) GetGridCopy() [][]Cell {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	copy := make([][]Cell, s.Grid.Size)
	for i := range copy {
		copy[i] = make([]Cell, s.Grid.Size)
		for j := range copy[i] {
			copy[i][j] = s.Grid.Cells[i][j]
		}
	}
	return copy
}

package simulation

import (
	"math/rand"
	"sync"
)

// StepFishParallel processes all fish in parallel using goroutines
func (s *Simulation) StepFishParallel(numThreads int) {
	fishList := make([]*Fish, 0)

	// Get current list of fish using Range (thread-safe iteration)
	s.Fish.Range(func(key, value interface{}) bool {
		fish := value.(*Fish)
		fishList = append(fishList, fish)
		return true
	})

	// Divide fish into chunks for parallel processing
	chunkSize := (len(fishList) + numThreads - 1) / numThreads
	var wg sync.WaitGroup

	for i := 0; i < numThreads && i*chunkSize < len(fishList); i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			start := threadID * chunkSize
			end := start + chunkSize
			if end > len(fishList) {
				end = len(fishList)
			}

			// Process this chunk of fish
			for idx := start; idx < end; idx++ {
				fish := fishList[idx]

				if _, exists := s.Fish.Load(fish.ID); !exists {
					continue
				}

				// Get adjacent empty cells
				empty, _, _ := s.Grid.GetAdjacentCells(fish.X, fish.Y)

				// Move fish if there's space
				if len(empty) > 0 {
					newPos := empty[rand.Intn(len(empty))]
					s.Grid.SetCell(fish.X, fish.Y, Cell{Type: EMPTY, ID: 0})
					fish.X = newPos.x
					fish.Y = newPos.y
					s.Grid.SetCell(fish.X, fish.Y, Cell{Type: FISH, ID: fish.ID})
				}

				// Increase age
				fish.IncreaseAge()

				// Check if fish can reproduce
				if fish.CanReproduce(s.FishBreedAge) {
					newFish := NewFish(s.NextFishID, fish.X, fish.Y, s.FishBreedAge)
					s.NextFishID++
					s.Fish.Store(newFish.ID, newFish)
					s.FishCount.Add(1)
					s.Grid.SetCell(newFish.X, newFish.Y, Cell{Type: FISH, ID: newFish.ID})
					fish.Age = 0
				}
			}
		}(i)
	}

	wg.Wait()
}

// StepSharksParallel processes all sharks in parallel using goroutines
func (s *Simulation) StepSharksParallel(numThreads int) {
	sharkList := make([]*Shark, 0)

	// Get current list of sharks using Range (thread-safe iteration)
	s.Sharks.Range(func(key, value interface{}) bool {
		shark := value.(*Shark)
		sharkList = append(sharkList, shark)
		return true
	})

	// Divide sharks into chunks for parallel processing
	chunkSize := (len(sharkList) + numThreads - 1) / numThreads
	var wg sync.WaitGroup

	for i := 0; i < numThreads && i*chunkSize < len(sharkList); i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			start := threadID * chunkSize
			end := start + chunkSize
			if end > len(sharkList) {
				end = len(sharkList)
			}

			// Process this chunk of sharks
			for idx := start; idx < end; idx++ {
				shark := sharkList[idx]

				if _, exists := s.Sharks.Load(shark.ID); !exists {
					continue
				}

				// Get adjacent cells
				empty, fish, _ := s.Grid.GetAdjacentCells(shark.X, shark.Y)

				// Sharks prioritize eating fish
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
		}(i)
	}

	wg.Wait()
}

// StepParallel advances simulation by one chronon using parallel processing
func (s *Simulation) StepParallel(numThreads int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if numThreads <= 1 {
		// Sequential mode
		s.stepFish()
		s.stepSharks()
	} else {
		// Parallel mode
		s.StepFishParallel(numThreads)
		s.StepSharksParallel(numThreads)
	}

	s.Chronon++
}

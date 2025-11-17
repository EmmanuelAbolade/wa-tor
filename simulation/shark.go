package simulation

// Shark represents a shark in the Wa-Tor world
type Shark struct {
	ID     int
	X      int
	Y      int
	Age    int
	Energy int
}

// NewShark creates a new shark at position (x, y)
func NewShark(id, x, y, initialEnergy int) *Shark {
	return &Shark{
		ID:     id,
		X:      x,
		Y:      y,
		Age:    0,
		Energy: initialEnergy,
	}
}

// CanReproduce checks if the shark is old enough to breed
func (s *Shark) CanReproduce(breedAge int) bool {
	return s.Age >= breedAge
}

// IsAlive checks if the shark still has energy
func (s *Shark) IsAlive() bool {
	return s.Energy > 0
}

// Starve reduces energy by 1 each chronon
func (s *Shark) Starve() {
	s.Energy--
}

// Eat increases energy by consuming a fish
func (s *Shark) Eat(energyPerFish int) {
	s.Energy += energyPerFish
}

// IncreaseAge increments the shark's age by 1
func (s *Shark) IncreaseAge() {
	s.Age++
}

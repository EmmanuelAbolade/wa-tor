package simulation

// Fish represents a fish in the Wa-Tor world
type Fish struct {
	ID  int
	X   int
	Y   int
	Age int
}

// NewFish creates a new fish at position (x, y)
func NewFish(id, x, y, breedAge int) *Fish {
	return &Fish{
		ID:  id,
		X:   x,
		Y:   y,
		Age: 0,
	}
}

// CanReproduce checks if the fish is old enough to breed
func (f *Fish) CanReproduce(breedAge int) bool {
	return f.Age >= breedAge
}

// IncreaseAge increments the fish's age by 1
func (f *Fish) IncreaseAge() {
	f.Age++
}

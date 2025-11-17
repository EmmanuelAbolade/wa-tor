package simulation

import (
	"sync"
)

// CellType represents what occupies a grid cell
type CellType int

const (
	EMPTY CellType = iota
	FISH
	SHARK
)

// Cell represents a single grid cell
type Cell struct {
	Type CellType
	ID   int
}

// Grid represents the toroidal world of Wa-Tor
type Grid struct {
	Size  int
	Cells [][]Cell
	mu    sync.RWMutex
}

// NewGrid creates and initializes a new grid
func NewGrid(size int) *Grid {
	grid := &Grid{
		Size:  size,
		Cells: make([][]Cell, size),
	}

	// Initialize empty cells
	for i := 0; i < size; i++ {
		grid.Cells[i] = make([]Cell, size)
		for j := 0; j < size; j++ {
			grid.Cells[i][j] = Cell{Type: EMPTY, ID: 0}
		}
	}

	return grid
}

// GetCell returns the cell at (x, y) with toroidal wrapping
func (g *Grid) GetCell(x, y int) Cell {
	g.mu.RLock()
	defer g.mu.RUnlock()

	x = (x + g.Size) % g.Size
	y = (y + g.Size) % g.Size

	return g.Cells[x][y]
}

// SetCell sets the cell at (x, y)
func (g *Grid) SetCell(x, y int, cell Cell) {
	g.mu.Lock()
	defer g.mu.Unlock()

	x = (x + g.Size) % g.Size
	y = (y + g.Size) % g.Size

	g.Cells[x][y] = cell
}

// GetAdjacentCells returns all 4 adjacent cells
func (g *Grid) GetAdjacentCells(x, y int) (empty []struct{ x, y int }, fish []struct{ x, y int }, sharks []struct{ x, y int }) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Define the 4 adjacent positions: North, South, East, West
	adjacentPositions := [][2]int{
		{x, y - 1}, // North
		{x, y + 1}, // South
		{x + 1, y}, // East
		{x - 1, y}, // West
	}

	for _, pos := range adjacentPositions {
		nx := (pos[0] + g.Size) % g.Size
		ny := (pos[1] + g.Size) % g.Size
		cell := g.Cells[nx][ny]

		switch cell.Type {
		case EMPTY:
			empty = append(empty, struct{ x, y int }{nx, ny})
		case FISH:
			fish = append(fish, struct{ x, y int }{nx, ny})
		case SHARK:
			sharks = append(sharks, struct{ x, y int }{nx, ny})
		}
	}

	return empty, fish, sharks
}

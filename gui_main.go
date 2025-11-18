package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abola/wa-tor/graphics"
	"github.com/abola/wa-tor/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"image/color"
)

// Game represents the simulation game
type Game struct {
	Simulation *simulation.Simulation
	Display    *graphics.VisualDisplay
	CellSize   int
	Paused     bool
}

// Update updates the game state
func (g *Game) Update() error {
	if !g.Paused {
		g.Simulation.Step()
	}
	return nil
}

// Draw draws the game
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	gridCopy := g.Simulation.GetGridCopy()
	g.Display.DrawGrid(screen, gridCopy)

	// Draw stats on screen
	chronon, numFish, numSharks := g.Simulation.GetStats()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Chronon: %d | Fish: %d | Sharks: %d", chronon, numFish, numSharks))
}

// Layout returns the screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.CellSize * g.Display.GridSize, g.CellSize * g.Display.GridSize
}

func main() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Simulation parameters
	gridSize := 50
	numFish := 500
	numSharks := 100
	fishBreedAge := 4
	sharkBreedAge := 8
	sharkStarveTime := 10
	cellSize := 10

	// Create simulation
	sim := simulation.NewSimulation(gridSize, numFish, numSharks, fishBreedAge, sharkBreedAge, sharkStarveTime)

	// Create visual display
	display := graphics.NewVisualDisplay(gridSize, cellSize)

	// Create game
	game := &Game{
		Simulation: sim,
		Display:    display,
		CellSize:   cellSize,
		Paused:     false,
	}

	// Set window title and size
	ebiten.SetWindowTitle("Wa-Tor Predator-Prey Simulation")
	ebiten.SetWindowResizable(false)

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error:", err)
	}
}

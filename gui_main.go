package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/abola/wa-tor/graphics"
	"github.com/abola/wa-tor/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	// Get grid directly without copying
	g.Simulation.Mu.RLock()
	gridSize := g.Simulation.Grid.Size
	gridCells := g.Simulation.Grid.Cells
	g.Simulation.Mu.RUnlock()

	// Draw grid cells directly
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			cell := gridCells[x][y]
			screenX := x * g.CellSize
			screenY := y * g.CellSize

			color := g.Display.ColorEmpty
			if cell.Type == 1 { // FISH
				color = g.Display.ColorFish
			} else if cell.Type == 2 { // SHARK
				color = g.Display.ColorShark
			}
			ebitenutil.DrawRect(screen, float64(screenX), float64(screenY), float64(g.CellSize), float64(g.CellSize), color)
		}
	}

	// Draw stats
	chronon, numFish, numSharks := g.Simulation.GetStats()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Chronon: %d | Fish: %d | Sharks: %d", chronon, numFish, numSharks))
}

// Layout returns the screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.CellSize * g.Display.GridSize, g.CellSize * g.Display.GridSize
}

// RunGUI starts the graphical interface
func RunGUI() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Simulation parameters
	gridSize := 50
	numFish := 300
	numSharks := 50
	fishBreedAge := 12
	sharkBreedAge := 12
	sharkStarveTime := 25
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

	// Set target FPS to keep simulation smooth
	ebiten.SetTPS(10) // 60 updates per second

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error:", err)
	}
}

package graphics

import (
	"github.com/abola/wa-tor/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// VisualDisplay manages the graphical window
type VisualDisplay struct {
	CellSize   int
	GridSize   int
	ColorEmpty color.Color
	ColorFish  color.Color
	ColorShark color.Color
}

// NewVisualDisplay creates a new visual display
func NewVisualDisplay(gridSize, cellSize int) *VisualDisplay {
	return &VisualDisplay{
		GridSize:   gridSize,
		CellSize:   cellSize,
		ColorEmpty: color.RGBA{100, 150, 255, 255}, // Light blue
		ColorFish:  color.RGBA{0, 255, 0, 255},     // Green
		ColorShark: color.RGBA{255, 0, 0, 255},     // Red
	}
}

// DrawGrid draws the grid on the screen
func (vd *VisualDisplay) DrawGrid(screen *ebiten.Image, grid [][]simulation.Cell) {
	for i := 0; i < vd.GridSize; i++ {
		for j := 0; j < vd.GridSize; j++ {
			cell := grid[i][j]
			x := float64(i * vd.CellSize)
			y := float64(j * vd.CellSize)

			// Determine color based on cell type
			var cellColor color.Color
			switch cell.Type {
			case simulation.EMPTY:
				cellColor = vd.ColorEmpty
			case simulation.FISH:
				cellColor = vd.ColorFish
			case simulation.SHARK:
				cellColor = vd.ColorShark
			}

			// Draw the cell as a rectangle
			ebitenutil.DrawRect(screen, x, y, float64(vd.CellSize), float64(vd.CellSize), cellColor)
		}
	}
}

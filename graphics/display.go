package graphics

import (
	"fmt"

	"github.com/abola/wa-tor/simulation"
)

// ANSI color codes
const (
	ColorReset = "\033[0m"
	ColorCyan  = "\033[36m" // Light blue for empty
	ColorGreen = "\033[32m" // Green for fish
	ColorRed   = "\033[31m" // Red for sharks
)

// DisplayGrid prints the grid with colors
func DisplayGrid(grid [][]simulation.Cell, gridSize int) {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			cell := grid[i][j]

			switch cell.Type {
			case simulation.EMPTY:
				fmt.Print(ColorCyan + "." + ColorReset)
			case simulation.FISH:
				fmt.Print(ColorGreen + "F" + ColorReset)
			case simulation.SHARK:
				fmt.Print(ColorRed + "S" + ColorReset)
			}
		}
		fmt.Println() // New line after each row
	}
	fmt.Println() // Extra line for spacing
}

# Wa-Tor: Predator-Prey Simulation

A concurrent simulation of sharks and fish on a toroidal grid, implemented in Go with real-time graphical visualization.

## Overview

Wa-Tor is an ecological simulation based on the Mathematical Recreations column by A.K. Dewdney. It models the interactions between two species:
- **Fish** (prey) that move randomly and reproduce
- **Sharks** (predators) that hunt fish, starve without food, and reproduce

The simulation runs on a toroidal (wrap-around) grid where both species follow simple behavioral rules.

## Features

- Real-time graphical visualization with Ebiten
- Toroidal grid (wraparound edges)
- Configurable parameters via command-line flags
- Population tracking and display
- Color-coded visualization:
    - Light blue = empty cells
    - Green = fish
    - Red = sharks

## Building

### Prerequisites
- Go 1.18 or later
- Ebiten dependencies (installed automatically)

### Build Instructions
```bash
git clone https://github.com/EmmanuelAbolade/wa-tor.git
cd wa-tor
go mod download
go build
```

## Running the Simulation

### Graphical Mode (default)
```bash
./wa-tor.exe
```

This opens a window showing the live simulation with default parameters:
- Grid Size: 50×50
- Fish: 500
- Sharks: 100
- Fish Breed Age: 4 chronons
- Shark Breed Age: 8 chronons
- Shark Starve Time: 10 chronons

## Parameters

| Parameter     | Default | Description                                            |
|---------------|---------|--------------------------------------------------------|
| `-grid`       | 50      | Grid dimensions (size × size)                          |
| `-fish`       | 500     | Initial fish population                                |
| `-sharks`     | 100     | Initial shark population                               |
| `-fishbreed`  | 4       | Chronons before fish can reproduce                     |
| `-sharkbreed` | 8       | Chronons before shark can reproduce                    |
| `-starve`     | 10      | Energy/chronons before shark starves                   |
| `-threads`    | 1       | Number of threads (for future parallel implementation) |

## Rules

### Fish
- Move randomly to adjacent unoccupied cells each chronon
- If all adjacent cells occupied, do not move
- Reproduce after reaching breed age, leaving offspring in old position
- Age resets after reproduction

### Sharks
- Hunt: Move to adjacent fish if available (prioritized)
- Otherwise: Move to random adjacent unoccupied cell
- Lose 1 energy each chronon
- Die if energy reaches 0
- Gain energy when eating fish (3 energy per fish)
- Reproduce after reaching breed age with split energy

## Project Structure
```
wa-tor/
├── main simulation logic
├── simulation/
│   ├── grid.go       - Toroidal grid implementation
│   ├── fish.go       - Fish entity and methods
│   ├── shark.go      - Shark entity and methods
│   └── simulation.go  - Core simulation engine
├── graphics/
│   ├── display.go    - ASCII console display
│   └── visual_display.go - Graphical Ebiten display
├── gui_main.go       - Entry point with graphical window
├── go.mod            - Go module definition
├── README.md         - This file
└── LICENSE           - Project license
```

## Technologies Used

- **Language:** Go
- **Graphics:** Ebiten v2 (2D game engine)
- **Concurrency:** Go routines and sync.RWMutex
- **Version Control:** Git & GitHub

## Performance Results

The simulation achieves concurrent speedup on multi-core systems:

| Threads | Time (ms) | Chronons/sec | Speedup |
|---------|-----------|--------------|---------|
| 1       | 83.67     | 597.62       | 1.00x   |
| 2       | 113.10    | 442.07       | 0.74x   |
| 4       | 137.12    | 364.64       | 0.61x   |
| 8       | 49.35     | 1,013.12     | 1.70x   |

See [PERFORMANCE_RESULTS.md](PERFORMANCE_RESULTS.md) for detailed analysis.

## Documentation

### Generating Doxygen Documentation

To generate API documentation from the source code comments:

**Prerequisites:**
- Doxygen must be installed

**Installation:**
- **Linux/Mac:** Usually pre-installed or `brew install doxygen`
- **Windows:** Download from https://www.doxygen.nl/download.html

**Generate Documentation:**
```bash
doxygen Doxyfile
```

Documentation will be generated in `docs/html/index.html`. Open this file in a web browser to view the API documentation.

**Project Doxygen Configuration:**
The `Doxyfile` is configured to extract:
- All Go source files in `simulation/` and `graphics/` packages
- Public and private functions/structs
- All inline documentation comments

## Future Enhancements

- [ ] Multithreaded simulation for parallel execution
- [ ] Parameter sliders for real-time adjustment
- [ ] Population graph/chart display
- [ ] Simulation pause/resume controls
- [ ] Performance benchmarking suite
- [x] Multithreaded simulation for parallel execution (1.70x speedup achieved)
## References

- Dewdney, A.K. (1984). Computer Recreations; Sharks and Fish wage an ecological war on the toroidal planet of Wa-Tor - Scientific American, pp. 14-22
- Wikipedia (2025). Wa-Tor Simulation. https://en.wikipedia.org/wiki/Wa-Tor accessed [30 October 2025].
- Leinweb (2025) Original Wa-Tor Description. https://www.leinweb.com/snackbar/wator/#:~:text=WATOR%20is%20a%20simulation%20of%20the%20interaction%20over,are%20stable%20when%20the%20area%20is%20made%20small. accessed [2 November 2025].

## Author

Emmanuel Abolade

## License

See LICENSE file for details
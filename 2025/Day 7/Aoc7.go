package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Position struct {
	row int
	col int
}

func countSplits(gridLines []string) int {
	// Normalize grid
	if len(gridLines) == 0 {
		return 0
	}

	// Convert to 2D grid
	grid := make([][]rune, len(gridLines))
	maxCols := 0
	for i, line := range gridLines {
		grid[i] = []rune(line)
		if len(grid[i]) > maxCols {
			maxCols = len(grid[i])
		}
	}

	R := len(grid)
	C := maxCols

	// Pad rows to equal width with spaces
	for i := 0; i < R; i++ {
		for len(grid[i]) < C {
			grid[i] = append(grid[i], ' ')
		}
	}

	// Find source 'S'
	var source *Position
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			if grid[r][c] == 'S' {
				source = &Position{r, c}
				break
			}
		}
		if source != nil {
			break
		}
	}

	if source == nil {
		panic("No source 'S' found in grid")
	}

	// Active beams as set of positions
	active := make(map[Position]bool)
	active[*source] = true
	splits := 0
	seenSplitPositions := make(map[Position]bool)

	for len(active) > 0 {
		newActive := make(map[Position]bool)

		// Process each beam: attempt to move one row down
		for pos := range active {
			r, c := pos.row, pos.col
			nr := r + 1

			if nr >= R {
				// Beam leaves the grid
				continue
			}

			cellBelow := grid[nr][c]
			if cellBelow == '^' {
				// Split occurs at (nr, c)
				splitPos := Position{nr, c}
				if !seenSplitPositions[splitPos] {
					splits++
					seenSplitPositions[splitPos] = true
				}

				// Spawn beams at immediate left and right of splitter (same row nr)
				leftCol := c - 1
				rightCol := c + 1

				if leftCol >= 0 && leftCol < C {
					newActive[Position{nr, leftCol}] = true
				}
				if rightCol >= 0 && rightCol < C {
					newActive[Position{nr, rightCol}] = true
				}
				// Original beam does not continue downward beyond the '^'
			} else {
				// Move into the cell below (covers '.' and 'S' or any char not '^')
				newActive[Position{nr, c}] = true
			}
		}

		// Continue with new active positions
		active = newActive
	}

	return splits
}

func countTimelines(gridLines []string) int {
	// Normalize grid
	if len(gridLines) == 0 {
		return 0
	}

	// Convert to 2D grid
	grid := make([][]rune, len(gridLines))
	maxCols := 0
	for i, line := range gridLines {
		grid[i] = []rune(line)
		if len(grid[i]) > maxCols {
			maxCols = len(grid[i])
		}
	}

	R := len(grid)
	C := maxCols

	// Pad rows to equal width with spaces
	for i := 0; i < R; i++ {
		for len(grid[i]) < C {
			grid[i] = append(grid[i], ' ')
		}
	}

	// Initialize beam strength matrix
	strength := make([][]int, R)
	for i := 0; i < R; i++ {
		strength[i] = make([]int, C)
	}

	// Find source 'S' and initialize strength
	activeCols := make(map[int]bool)
	for c := 0; c < C; c++ {
		if grid[0][c] == 'S' {
			strength[0][c] = 1
			activeCols[c] = true
		}
	}

	// Process row by row
	for y := 0; y < R-1; y++ {
		nextActiveCols := make(map[int]bool)

		for x := range activeCols {
			// Get cell at next row
			cell := grid[y+1][x]

			// Check if cell above is a splitter
			cellAbove := grid[y][x]
			isBelowSplitter := (cellAbove == '^')

			if cell == '^' {
				// This row has a splitter: split the beam left and right
				left := x - 1
				right := x + 1

				if left >= 0 && left < C {
					strength[y+1][left] += strength[y][x]
					nextActiveCols[left] = true
				}
				if right >= 0 && right < C {
					strength[y+1][right] += strength[y][x]
					nextActiveCols[right] = true
				}
			} else {
				// Regular cell or 'S': beam continues if not directly below a splitter
				if !isBelowSplitter {
					strength[y+1][x] += strength[y][x]
					nextActiveCols[x] = true
				}
			}
		}

		activeCols = nextActiveCols
	}

	// Return sum of strengths at the bottom row (beams reaching exit)
	total := 0
	for c := 0; c < C; c++ {
		total += strength[R-1][c]
	}

	return total
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_7")

	// Read file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Handle both Windows (\r\n) and Unix (\n) line endings
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	// Remove empty lines at the end
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Part 1
	result1 := countSplits(lines)
	fmt.Printf("Day 7 - Part 1 (Total splits): %d\n", result1)

	// Part 2
	result2 := countTimelines(lines)
	fmt.Printf("Day 7 - Part 2 (Number of timelines): %d\n", result2)
}

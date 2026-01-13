package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Day 4 — Warehouse Rolls
// Part 1: Count accessible '@' rolls (rolls with fewer than 4 adjacent '@' neighbors)
// Part 2: Iteratively remove accessible rolls until none remain, count total removed

// countAccessibleRolls counts rolls that have fewer than 4 adjacent '@' neighbors
func countAccessibleRolls(grid [][]rune) int {
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)
	cols := len(grid[0])

	// 8 direction offsets: up-left, up, up-right, left, right, down-left, down, down-right
	neighbors := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	accessibleCount := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '@' {
				continue
			}

			// Count adjacent '@' rolls
			adjAt := 0
			for _, offset := range neighbors {
				nr, nc := r+offset[0], c+offset[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
					adjAt++
					if adjAt >= 4 {
						break // Early exit optimization
					}
				}
			}

			// Accessible if fewer than 4 neighbors
			if adjAt < 4 {
				accessibleCount++
			}
		}
	}

	return accessibleCount
}

// findAccessiblePositions returns list of all accessible roll positions
func findAccessiblePositions(grid [][]rune) [][2]int {
	if len(grid) == 0 {
		return nil
	}

	rows := len(grid)
	cols := len(grid[0])

	neighbors := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	accessible := [][2]int{}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '@' {
				continue
			}

			adjAt := 0
			for _, offset := range neighbors {
				nr, nc := r+offset[0], c+offset[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
					adjAt++
					if adjAt >= 4 {
						break
					}
				}
			}

			if adjAt < 4 {
				accessible = append(accessible, [2]int{r, c})
			}
		}
	}

	return accessible
}

// removeIteratively removes accessible rolls until none remain
func removeIteratively(grid [][]rune, verbose bool) int {
	totalRemoved := 0
	roundNo := 0

	for {
		accessible := findAccessiblePositions(grid)
		if len(accessible) == 0 {
			break
		}

		roundNo++
		if verbose {
			fmt.Printf("Round %d: removing %d rolls\n", roundNo, len(accessible))
		}

		// Remove all accessible rolls simultaneously
		for _, pos := range accessible {
			grid[pos[0]][pos[1]] = '.'
		}

		totalRemoved += len(accessible)
	}

	return totalRemoved
}

// readGrid reads grid from file
func readGrid(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	grid := [][]rune{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		if line == "" {
			continue
		}
		grid = append(grid, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return grid, nil
}

// copyGrid creates a deep copy of the grid
func copyGrid(grid [][]rune) [][]rune {
	copy := make([][]rune, len(grid))
	for i := range grid {
		copy[i] = make([]rune, len(grid[i]))
		for j := range grid[i] {
			copy[i][j] = grid[i][j]
		}
	}
	return copy
}

func solve(filename string) (int, int, error) {
	grid, err := readGrid(filename)
	if err != nil {
		return 0, 0, err
	}

	// Part 1: Count initial accessible rolls
	part1 := countAccessibleRolls(grid)

	// Part 2: Iteratively remove rolls (need a copy since we modify the grid)
	gridCopy := copyGrid(grid)
	part2 := removeIteratively(gridCopy, false)

	return part1, part2, nil
}

func main() {
	// Default input file is `input_day_4` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "input_day_4")

	inputFile := defaultFile
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", inputFile)
		os.Exit(2)
	}

	part1, part2, err := solve(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("PART 1: Initial accessible rolls")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Accessible rolls: %d\n", part1)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("PART 2: Remove rolls iteratively")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total rolls removed: %d\n", part2)
}

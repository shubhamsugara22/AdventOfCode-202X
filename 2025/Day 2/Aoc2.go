package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Day 2 — Invalid IDs
// Part 1: Sum IDs where first half of digits equals second half (even digit count)
// Part 2: Sum IDs that can be represented as a repeating block pattern

// isInvalidID checks if a number has even digit count and first half equals second half
func isInvalidID(n int) bool {
	s := strconv.Itoa(n)
	// Must have even number of digits
	if len(s)%2 != 0 {
		return false
	}
	mid := len(s) / 2
	// First half repeated twice
	return s[:mid] == s[mid:]
}

// isInvalidIDPart2 checks if a number can be represented as a repeating block
func isInvalidIDPart2(n int) bool {
	s := strconv.Itoa(n)
	L := len(s)

	// Try all possible block lengths
	for k := 1; k <= L/2; k++ {
		if L%k != 0 {
			continue // must divide evenly
		}

		repeatCount := L / k
		if repeatCount < 2 {
			continue // must repeat at least twice
		}

		block := s[:k]

		// Check if repeating the block gives us the original string
		repeated := strings.Repeat(block, repeatCount)
		if repeated == s {
			return true
		}
	}

	return false
}

func solve(filename string) (int, int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, 0, err
	}

	line := strings.TrimSpace(string(data))
	ranges := strings.Split(line, ",")

	totalPart1 := 0
	totalPart2 := 0

	for _, r := range ranges {
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			continue
		}

		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}

		for n := start; n <= end; n++ {
			if isInvalidID(n) {
				totalPart1 += n
			}
			if isInvalidIDPart2(n) {
				totalPart2 += n
			}
		}
	}

	return totalPart1, totalPart2, nil
}

func main() {
	// Default input file is `input_day_2` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "input_day_2")

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

	fmt.Printf("Sum of all invalid IDs = %d\n", part1)
	fmt.Printf("Part 2 sum of invalid IDs = %d\n", part2)
}

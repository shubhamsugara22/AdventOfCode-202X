package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Day 1 — Secret Entrance
// Part 1: count times the dial is at 0 after a rotation finishes
// Part 2: count times the dial is at 0 during any click while performing rotations
//
// Expects input file at Input_day_1 (one instruction per line like "L68" or "R48").

func solve(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// initial dial position
	pos := 50

	// Part 1: count times position == 0 after a rotation completes
	part1Count := 0

	// Part 2: count times position == 0 during any click while performing rotations
	part2Count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if len(line) < 2 {
			// skip malformed lines
			continue
		}

		direction := strings.ToUpper(string(line[0]))
		stepsStr := line[1:]

		steps, err := strconv.Atoi(stepsStr)
		if err != nil {
			// skip malformed lines
			continue
		}

		if direction != "R" && direction != "L" {
			// ignore unknown directions
			continue
		}

		// --- Part 2: count intermediate hits of 0 during this rotation ---
		// For right (increasing): position at click k is (pos + k) % 100
		// For left (decreasing): position at click k is (pos - k) % 100
		// We want number of integers k in [1..steps] such that the expression == 0.
		var k0 int
		if direction == "R" {
			// solve (pos + k) % 100 == 0 -> k ≡ (100 - pos) (mod 100)
			k0 = (100 - pos) % 100
		} else { // "L"
			// solve (pos - k) % 100 == 0 -> k ≡ pos (mod 100)
			k0 = pos % 100
		}

		// convert k0 of 0 to 100 because k=0 is not a click;
		// the first click that hits 0 occurs at 100 clicks when k0==0
		if k0 == 0 {
			k0 = 100
		}

		// if the first k0 occurs within the number of steps,
		// we get 1 + floor((steps - k0)/100) total hits
		if steps >= k0 {
			hits := 1 + (steps-k0)/100
			part2Count += hits
		}

		// --- Now update the position (end of rotation) for both parts ---
		if direction == "R" {
			pos = (pos + steps) % 100
		} else {
			pos = (pos - steps) % 100
		}

		// ensure position is non-negative
		if pos < 0 {
			pos += 100
		}

		// Part 1: if after this rotation the dial is at 0, count it
		if pos == 0 {
			part1Count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return part1Count, part2Count, nil
}

func main() {
	// Default input file is `Input_day_1` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "Input_day_1")

	inputFile := defaultFile
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", inputFile)
		os.Exit(2)
	}

	p1, p2, err := solve(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1 (count ends at 0): %d\n", p1)
	fmt.Printf("Part 2 (count any click at 0): %d\n", p2)
}

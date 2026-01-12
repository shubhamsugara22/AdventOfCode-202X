package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Day 3 — Joltage Banks
// Part 1: For each bank, find max 2-digit number from digits in order (sum all banks)
// Part 2: For each bank, find max 12-digit number using monotonic stack (sum all banks)

// maxKDigits returns the largest possible number formed by keeping exactly k digits
// in the same order using a monotonic stack algorithm
func maxKDigits(numStr string, k int) string {
	remove := len(numStr) - k
	stack := []rune{}

	for _, ch := range numStr {
		for remove > 0 && len(stack) > 0 && stack[len(stack)-1] < ch {
			stack = stack[:len(stack)-1]
			remove--
		}
		stack = append(stack, ch)
	}

	return string(stack[:k])
}

func solve(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	totalPart1 := 0 // For k = 2
	totalPart2 := 0 // For k = 12

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}

		// Part 1: choose best 2 digits
		best2Str := maxKDigits(s, 2)
		best2, err := strconv.Atoi(best2Str)
		if err != nil {
			continue
		}
		totalPart1 += best2

		// Part 2: choose best 12 digits
		best12Str := maxKDigits(s, 12)
		best12, err := strconv.Atoi(best12Str)
		if err != nil {
			continue
		}
		totalPart2 += best12
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return totalPart1, totalPart2, nil
}

func main() {
	// Default input file is `input_day_3` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "input_day_3")

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

	fmt.Printf("Part 1 Total Joltage: %d\n", part1)
	fmt.Printf("Part 2 Total Joltage: %d\n", part2)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Day 5 — Fresh Ingredient IDs
// Part 1: Count IDs that fall within any safe range
// Part 2: Count total IDs covered by merged safe ranges

// countFreshIDs counts IDs that fall within the safe ranges
func countFreshIDs(inputText string) int {
	lines := strings.Split(strings.TrimSpace(inputText), "\n")

	// Find blank line separating ranges from IDs
	blankIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			blankIndex = i
			break
		}
	}

	if blankIndex == -1 {
		return 0
	}

	rangeLines := lines[:blankIndex]
	idLines := lines[blankIndex+1:]

	// Parse ranges as [start, end] pairs
	type Range struct {
		start, end int
	}
	ranges := []Range{}

	for _, line := range rangeLines {
		parts := strings.Split(strings.TrimSpace(line), "-")
		if len(parts) != 2 {
			continue
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			ranges = append(ranges, Range{start, end})
		}
	}

	// Parse ingredient IDs
	ids := []int{}
	for _, line := range idLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err == nil {
			ids = append(ids, val)
		}
	}

	// Count fresh ones (IDs within any safe range)
	freshCount := 0
	for _, val := range ids {
		for _, r := range ranges {
			if val >= r.start && val <= r.end {
				freshCount++
				break // no need to check other ranges
			}
		}
	}

	return freshCount
}

// countFreshIDsPart2 counts total IDs covered by merged safe ranges
func countFreshIDsPart2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Find blank line
	blankIndex := -1
	for i, line := range lines {
		if line == "" {
			blankIndex = i
			break
		}
	}

	if blankIndex == -1 {
		return 0, nil
	}

	freshRanges := lines[:blankIndex]

	// Parse intervals
	type Interval struct {
		start, end int
	}
	intervals := []Interval{}

	for _, r := range freshRanges {
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			continue
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			intervals = append(intervals, Interval{start, end})
		}
	}

	// Sort ranges by start
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	if len(intervals) == 0 {
		return 0, nil
	}

	// Merge the intervals
	merged := []Interval{}
	currentStart := intervals[0].start
	currentEnd := intervals[0].end

	for _, interval := range intervals[1:] {
		s, e := interval.start, interval.end
		if s <= currentEnd+1 { // overlapping or touching
			if e > currentEnd {
				currentEnd = e
			}
		} else {
			merged = append(merged, Interval{currentStart, currentEnd})
			currentStart = s
			currentEnd = e
		}
	}
	merged = append(merged, Interval{currentStart, currentEnd})

	// Count total IDs covered by merged ranges
	totalFreshIDs := 0
	for _, interval := range merged {
		totalFreshIDs += interval.end - interval.start + 1
	}

	return totalFreshIDs, nil
}

func solve(filename string) (int, int, error) {
	// Read entire file for Part 1
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, 0, err
	}

	part1 := countFreshIDs(string(data))
	part2, err := countFreshIDsPart2(filename)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func main() {
	// Default input file is `input_day_5` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "input_day_5")

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
	fmt.Println("PART 1: Count fresh IDs by checking each ID")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Fresh IDs count: %d\n", part1)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("PART 2: Count fresh IDs by merging ranges")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total fresh IDs (by merged ranges): %d\n", part2)
}

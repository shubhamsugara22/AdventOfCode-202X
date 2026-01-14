package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Day 6 — Cephalopod Math Worksheet
// Part 1: Read vertical numbers with operators, calculate sum of all problems
// Part 2: Numbers written right-to-left in columns with digits stacked vertically

// isBlankColumn checks if a column is entirely blank/whitespace
func isBlankColumn(lines []string, col int) bool {
	for _, line := range lines {
		if col < len(line) && !isSpace(rune(line[col])) {
			return false
		}
	}
	return true
}

// isSpace checks if a rune is whitespace
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// findProblems identifies problem column ranges by finding blank column separators
func findProblems(lines []string, maxLen int) [][2]int {
	numRows := len(lines)
	problems := [][2]int{}
	inProblem := false
	start := 0

	for col := 0; col < maxLen; col++ {
		isBlank := true
		for row := 0; row < numRows; row++ {
			if col < len(lines[row]) && !isSpace(rune(lines[row][col])) {
				isBlank = false
				break
			}
		}

		if !isBlank && !inProblem {
			inProblem = true
			start = col
		} else if isBlank && inProblem {
			inProblem = false
			problems = append(problems, [2]int{start, col})
		}
	}

	if inProblem {
		problems = append(problems, [2]int{start, maxLen})
	}

	return problems
}

// solveDay6 solves Part 1: vertical numbers with operators
func solveDay6(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	// Handle both Windows (\r\n) and Unix (\n) line endings
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	// Remove empty lines at the end
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return 0, nil
	}

	// Normalize line lengths
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Pad lines to maxLen
	for i := range lines {
		if len(lines[i]) < maxLen {
			lines[i] = lines[i] + strings.Repeat(" ", maxLen-len(lines[i]))
		}
	}

	numRows := len(lines)
	problems := findProblems(lines, maxLen)

	total := 0
	for _, prob := range problems {
		start, end := prob[0], prob[1]

		// Collect numbers from all rows except last
		numbers := []int{}
		for row := 0; row < numRows-1; row++ {
			chunk := strings.TrimSpace(lines[row][start:end])
			if chunk == "" {
				continue
			}
			// Take last token (right-aligned numbers)
			tokens := strings.Fields(chunk)
			if len(tokens) > 0 {
				num, err := strconv.Atoi(tokens[len(tokens)-1])
				if err == nil {
					numbers = append(numbers, num)
				}
			}
		}

		// Operator is on the last row
		opChunk := strings.TrimSpace(lines[numRows-1][start:end])
		op := rune(0)
		for _, ch := range opChunk {
			if ch == '+' || ch == '*' {
				op = ch
				break
			}
		}

		if len(numbers) == 0 || op == 0 {
			continue
		}

		// Compute the problem result
		res := numbers[0]
		if op == '+' {
			for _, n := range numbers[1:] {
				res += n
			}
		} else { // '*'
			for _, n := range numbers[1:] {
				res *= n
			}
		}

		total += res
	}

	return total, nil
}

// solveDay6Part2 solves Part 2: numbers written right-to-left in columns
func solveDay6Part2(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	// Handle both Windows (\r\n) and Unix (\n) line endings
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	// Remove empty lines at the end
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return 0, nil
	}

	// Normalize line lengths
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Pad lines to maxLen
	for i := range lines {
		if len(lines[i]) < maxLen {
			lines[i] = lines[i] + strings.Repeat(" ", maxLen-len(lines[i]))
		}
	}

	numRows := len(lines)
	problems := findProblems(lines, maxLen)

	total := 0
	for _, prob := range problems {
		start, end := prob[0], prob[1]

		// Operator for this problem
		opChunk := lines[numRows-1][start:end]
		op := rune(0)
		for _, ch := range opChunk {
			if ch == '+' || ch == '*' {
				op = ch
				break
			}
		}

		if op == 0 {
			continue
		}

		// Collect numbers by reading columns right-to-left
		numbers := []int{}
		for col := end - 1; col >= start; col-- {
			// Build digit string top->bottom from rows 0..numRows-2
			digits := []rune{}
			for row := 0; row < numRows-1; row++ {
				ch := rune(lines[row][col])
				if !isSpace(ch) {
					digits = append(digits, ch)
				}
			}

			if len(digits) == 0 {
				continue
			}

			numStr := string(digits)
			num, err := strconv.Atoi(numStr)
			if err == nil {
				numbers = append(numbers, num)
			}
		}

		if len(numbers) == 0 {
			continue
		}

		// Compute result
		res := numbers[0]
		if op == '+' {
			for _, n := range numbers[1:] {
				res += n
			}
		} else {
			for _, n := range numbers[1:] {
				res *= n
			}
		}

		total += res
	}

	return total, nil
}

func main() {
	// Default input file is `input_day_6` in the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	defaultFile := filepath.Join(cwd, "input_day_6")

	inputFile := defaultFile
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", inputFile)
		os.Exit(2)
	}

	part1, err := solveDay6(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving part 1: %v\n", err)
		os.Exit(1)
	}

	part2, err := solveDay6Part2(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving part 2: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Day 6 - Part 1: %d\n", part1)
	fmt.Printf("Day 6 - Part 2: %d\n", part2)
}

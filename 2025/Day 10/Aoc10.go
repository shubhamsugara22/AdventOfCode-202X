package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func parseLine(line string) (string, [][]int, error) {
	// Pattern: [lights](buttons...){...}
	pattern := regexp.MustCompile(`\[([.#]+)\](.*)`)
	matches := pattern.FindStringSubmatch(line)
	if len(matches) < 3 {
		return "", nil, fmt.Errorf("invalid line format")
	}

	lights := matches[1]
	rest := matches[2]

	// Extract buttons from parentheses
	buttonPattern := regexp.MustCompile(`\((.*?)\)`)
	buttonMatches := buttonPattern.FindAllStringSubmatch(rest, -1)

	var buttons [][]int
	for _, match := range buttonMatches {
		if len(match) < 2 {
			continue
		}
		content := strings.TrimSpace(match[1])
		if content == "" {
			buttons = append(buttons, []int{})
		} else {
			parts := strings.Split(content, ",")
			var btn []int
			for _, p := range parts {
				num, err := strconv.Atoi(strings.TrimSpace(p))
				if err == nil {
					btn = append(btn, num)
				}
			}
			buttons = append(buttons, btn)
		}
	}

	return lights, buttons, nil
}

func solveMachine(lights string, buttons [][]int) int {
	n := len(lights)
	m := len(buttons)

	// Target state
	target := make([]int, n)
	for i, c := range lights {
		if c == '#' {
			target[i] = 1
		}
	}

	// Build matrix B: size n x m
	B := make([][]int, n)
	for i := range B {
		B[i] = make([]int, m)
	}

	for j, btn := range buttons {
		for _, bit := range btn {
			if bit < n {
				B[bit][j] ^= 1
			}
		}
	}

	// Gaussian elimination
	row := 0
	where := make([]int, n)
	for i := range where {
		where[i] = -1
	}

	for col := 0; col < m; col++ {
		// Find pivot
		sel := -1
		for r := row; r < n; r++ {
			if B[r][col] == 1 {
				sel = r
				break
			}
		}
		if sel == -1 {
			continue
		}

		// Swap rows
		B[row], B[sel] = B[sel], B[row]
		target[row], target[sel] = target[sel], target[row]
		where[row] = col

		// Eliminate
		for r := 0; r < n; r++ {
			if r != row && B[r][col] == 1 {
				for c := col; c < m; c++ {
					B[r][c] ^= B[row][c]
				}
				target[r] ^= target[row]
			}
		}

		row++
	}

	// Check for inconsistency
	for r := row; r < n; r++ {
		if target[r] == 1 {
			return -1 // no solution
		}
	}

	// Find free variables
	pivotCols := make(map[int]bool)
	for r := 0; r < n; r++ {
		if where[r] != -1 {
			pivotCols[where[r]] = true
		}
	}

	var freeVars []int
	for col := 0; col < m; col++ {
		if !pivotCols[col] {
			freeVars = append(freeVars, col)
		}
	}

	// Try all free variable combinations
	best := -1
	numFree := len(freeVars)
	totalCombos := 1 << uint(numFree)

	for combo := 0; combo < totalCombos; combo++ {
		x := make([]int, m)

		// Set free variables
		for idx, col := range freeVars {
			if (combo>>uint(idx))&1 == 1 {
				x[col] = 1
			}
		}

		// Determine pivot variable values
		for r := 0; r < n; r++ {
			col := where[r]
			if col != -1 {
				val := target[r]
				for c := col + 1; c < m; c++ {
					if B[r][c] == 1 {
						val ^= x[c]
					}
				}
				x[col] = val
			}
		}

		// Count weight
		weight := 0
		for _, v := range x {
			weight += v
		}

		if best == -1 || weight < best {
			best = weight
		}
	}

	return best
}

func solveDay10Part1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lights, buttons, err := parseLine(line)
		if err != nil {
			continue
		}

		presses := solveMachine(lights, buttons)
		if presses >= 0 {
			total += presses
		}
	}

	return total, scanner.Err()
}

// Part 2: Linear Programming approach
// Note: Go doesn't have a built-in LP solver like Python's PuLP.
// For a complete Part 2 solution, you would need to use an external LP library
// such as github.com/draffensperger/golp or implement a simplex solver.
//
// The problem formulation would be:
// - Minimize: sum of all button press variables
// - Subject to: for each voltage/joltage position i,
//               sum of button_j variables where i is in buttons[j] == voltage[i]
//
// Below is a placeholder that parses the input format for Part 2.

func parseLine2(line string) ([][]int, []int, error) {
	// Extract buttons from parentheses
	buttonPattern := regexp.MustCompile(`\((.*?)\)`)
	buttonMatches := buttonPattern.FindAllStringSubmatch(line, -1)

	var buttons [][]int
	for _, match := range buttonMatches {
		if len(match) < 2 {
			continue
		}
		content := strings.TrimSpace(match[1])
		if content == "" {
			buttons = append(buttons, []int{})
		} else {
			parts := strings.Split(content, ",")
			var btn []int
			for _, p := range parts {
				num, err := strconv.Atoi(strings.TrimSpace(p))
				if err == nil {
					btn = append(btn, num)
				}
			}
			buttons = append(buttons, btn)
		}
	}

	// Extract voltages from braces
	voltagePattern := regexp.MustCompile(`\{(.*?)\}`)
	voltageMatch := voltagePattern.FindStringSubmatch(line)
	var voltages []int
	if len(voltageMatch) >= 2 {
		parts := strings.Split(voltageMatch[1], ",")
		for _, p := range parts {
			num, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				voltages = append(voltages, num)
			}
		}
	}

	return buttons, voltages, nil
}

func solveDay10Part2(filename string) {
	fmt.Println("Part 2 requires an Integer Linear Programming solver.")
	fmt.Println("Consider using an external library like:")
	fmt.Println("  - github.com/draffensperger/golp")
	fmt.Println("  - Or implement a simplex/branch-and-bound solver")
	fmt.Println("\nProblem formulation:")
	fmt.Println("  Minimize: sum of button presses")
	fmt.Println("  Subject to: voltage constraints at each position")
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_10")

	result, err := solveDay10Part1(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Day 10 Part 1: %d\n", result)
	fmt.Println()
	solveDay10Part2(inputPath)
}

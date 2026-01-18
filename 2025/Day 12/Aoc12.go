package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Region struct {
	width  int
	height int
	counts []int
}

func parseInput(filename string) (map[int][]string, []Region, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	shapes := make(map[int][]string)
	var regions []Region
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r\n")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	i := 0
	// Parse shapes
	for i < len(lines) {
		if lines[i] == "" {
			i++
			continue
		}
		if !strings.Contains(lines[i], ":") {
			break
		}
		left := strings.Split(lines[i], ":")[0]
		idx, err := strconv.Atoi(left)
		if err != nil {
			break
		}
		i++
		var grid []string
		for i < len(lines) && lines[i] != "" {
			grid = append(grid, lines[i])
			i++
		}
		shapes[idx] = grid
		i++
	}

	// Parse regions
	for i < len(lines) {
		if lines[i] == "" {
			i++
			continue
		}
		parts := strings.Split(lines[i], ":")
		if len(parts) != 2 {
			i++
			continue
		}
		prefix := parts[0]
		nums := parts[1]

		dims := strings.Split(prefix, "x")
		if len(dims) != 2 {
			i++
			continue
		}
		W, _ := strconv.Atoi(dims[0])
		H, _ := strconv.Atoi(dims[1])

		countStrs := strings.Fields(nums)
		var counts []int
		for _, s := range countStrs {
			n, _ := strconv.Atoi(s)
			counts = append(counts, n)
		}

		regions = append(regions, Region{
			width:  W,
			height: H,
			counts: counts,
		})
		i++
	}

	return shapes, regions, nil
}

func solveDay12Part1(filename string) int {
	shapes, regions, err := parseInput(filename)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return 0
	}

	// Compute area of each shape
	shapeArea := make(map[int]int)
	for idx, grid := range shapes {
		area := 0
		for _, row := range grid {
			area += strings.Count(row, "#")
		}
		shapeArea[idx] = area
	}

	good := 0
	for _, region := range regions {
		required := 0
		for i := 0; i < len(region.counts); i++ {
			required += shapeArea[i] * region.counts[i]
		}
		available := region.width * region.height
		if required <= available {
			good++
		}
	}

	return good
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_12")

	fmt.Println("Day 12 - Part 1:", solveDay12Part1(inputPath))
}

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

type Point struct {
	x, y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readInput(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		points = append(points, Point{x, y})
	}
	return points, scanner.Err()
}

func solveDay9Part1(points []Point) int {
	n := len(points)
	maxArea := 0

	// Compare all pairs of points
	for i := 0; i < n; i++ {
		x1, y1 := points[i].x, points[i].y
		for j := i + 1; j < n; j++ {
			x2, y2 := points[j].x, points[j].y

			// Must form a valid rectangle (different x and y)
			if x1 != x2 && y1 != y2 {
				width := abs(x1-x2) + 1
				height := abs(y1-y2) + 1
				area := width * height
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea
}

// Geometry helpers for Part 2
func orient(a, b, c Point) int {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func segmentsProperIntersect(a, b, c, d Point) bool {
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)

	if (o1 > 0 && o2 < 0 || o1 < 0 && o2 > 0) && (o3 > 0 && o4 < 0 || o3 < 0 && o4 > 0) {
		return true
	}
	return false
}

func pointOnSegment(px, py, ax, ay, bx, by int) bool {
	// Integer collinearity and within bounding box
	if (bx-ax)*(py-ay) != (by-ay)*(px-ax) {
		return false
	}
	return min(ax, bx) <= px && px <= max(ax, bx) && min(ay, by) <= py && py <= max(ay, by)
}

func pointInPolygon(x, y int, poly []Point) bool {
	n := len(poly)

	// On-edge test
	for i := 0; i < n; i++ {
		x1, y1 := poly[i].x, poly[i].y
		x2, y2 := poly[(i+1)%n].x, poly[(i+1)%n].y
		if pointOnSegment(x, y, x1, y1, x2, y2) {
			return true
		}
	}

	// Ray casting
	inside := false
	for i := 0; i < n; i++ {
		x1, y1 := poly[i].x, poly[i].y
		x2, y2 := poly[(i+1)%n].x, poly[(i+1)%n].y

		if (y1 > y) != (y2 > y) {
			xinters := float64(x1) + float64(y-y1)*float64(x2-x1)/float64(y2-y1)
			if xinters >= float64(x) {
				inside = !inside
			}
		}
	}
	return inside
}

func rectangleValid(p1, p2 Point, polygon []Point, bboxMinX, bboxMaxX, bboxMinY, bboxMaxY int) bool {
	x1, y1 := p1.x, p1.y
	x2, y2 := p2.x, p2.y

	if x1 == x2 || y1 == y2 {
		return false
	}

	minX, maxX := min(x1, x2), max(x1, x2)
	minY, maxY := min(y1, y2), max(y1, y2)

	// Quick bbox reject
	if minX < bboxMinX || maxX > bboxMaxX || minY < bboxMinY || maxY > bboxMaxY {
		return false
	}

	// All rectangle corners must be inside or on polygon boundary
	corners := []Point{{x1, y1}, {x1, y2}, {x2, y1}, {x2, y2}}
	for _, c := range corners {
		if !pointInPolygon(c.x, c.y, polygon) {
			return false
		}
	}

	// Check no rectangle edge properly crosses polygon edges
	rectEdges := [][2]Point{
		{corners[0], corners[1]},
		{corners[1], corners[3]},
		{corners[3], corners[2]},
		{corners[2], corners[0]},
	}

	n := len(polygon)
	for _, edge := range rectEdges {
		a, b := edge[0], edge[1]
		axMin, axMax := min(a.x, b.x), max(a.x, b.x)
		ayMin, ayMax := min(a.y, b.y), max(a.y, b.y)

		for i := 0; i < n; i++ {
			c := polygon[i]
			d := polygon[(i+1)%n]

			// Quick bbox overlap test
			if axMax < min(c.x, d.x) || axMin > max(c.x, d.x) ||
				ayMax < min(c.y, d.y) || ayMin > max(c.y, d.y) {
				continue
			}

			if segmentsProperIntersect(a, b, c, d) {
				return false
			}
		}
	}

	return true
}

func solveDay9Part2(points []Point) int {
	if len(points) == 0 {
		return 0
	}

	polygon := points
	n := len(points)

	// Calculate bounding box
	bboxMinX, bboxMaxX := polygon[0].x, polygon[0].x
	bboxMinY, bboxMaxY := polygon[0].y, polygon[0].y
	for _, p := range polygon {
		bboxMinX = min(bboxMinX, p.x)
		bboxMaxX = max(bboxMaxX, p.x)
		bboxMinY = min(bboxMinY, p.y)
		bboxMaxY = max(bboxMaxY, p.y)
	}

	best := 0
	var bestPair [2]Point

	// Sort indices by x then y for better pruning
	idxs := make([]int, n)
	for i := 0; i < n; i++ {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i, j int) bool {
		if points[idxs[i]].x != points[idxs[j]].x {
			return points[idxs[i]].x < points[idxs[j]].x
		}
		return points[idxs[i]].y < points[idxs[j]].y
	})

	// Iterate pairs with pruning
	for aI := 0; aI < n; aI++ {
		i := idxs[aI]
		x1, y1 := points[i].x, points[i].y

		for bI := aI + 1; bI < n; bI++ {
			j := idxs[bI]
			x2, y2 := points[j].x, points[j].y

			if x1 == x2 || y1 == y2 {
				continue
			}

			potentialArea := (abs(x1-x2) + 1) * (abs(y1-y2) + 1)
			if potentialArea <= best {
				continue
			}

			if rectangleValid(points[i], points[j], polygon, bboxMinX, bboxMaxX, bboxMinY, bboxMaxY) {
				best = potentialArea
				bestPair = [2]Point{points[i], points[j]}
			}
		}
	}

	fmt.Printf("Best corners: %v\n", bestPair)
	return best
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_9")

	points, err := readInput(inputPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Println("Answer part 1:", solveDay9Part1(points))
	fmt.Println("Part 2 largest rectangle area (red+green):", solveDay9Part2(points))
}

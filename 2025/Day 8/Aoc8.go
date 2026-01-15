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

// Point represents a 3D point
type Point struct {
	x, y, z int
}

// DSU (Disjoint Set Union / Union-Find)
type DSU struct {
	p  []int
	r  []int
	sz []int
}

func NewDSU(n int) *DSU {
	dsu := &DSU{
		p:  make([]int, n),
		r:  make([]int, n),
		sz: make([]int, n),
	}
	for i := 0; i < n; i++ {
		dsu.p[i] = i
		dsu.sz[i] = 1
	}
	return dsu
}

func (dsu *DSU) find(a int) int {
	// Path compression
	for dsu.p[a] != a {
		dsu.p[a] = dsu.p[dsu.p[a]]
		a = dsu.p[a]
	}
	return a
}

func (dsu *DSU) union(a, b int) bool {
	ra := dsu.find(a)
	rb := dsu.find(b)
	if ra == rb {
		return false
	}
	if dsu.r[ra] < dsu.r[rb] {
		ra, rb = rb, ra
	}
	dsu.p[rb] = ra
	dsu.sz[ra] += dsu.sz[rb]
	if dsu.r[ra] == dsu.r[rb] {
		dsu.r[ra]++
	}
	return true
}

func (dsu *DSU) size(a int) int {
	return dsu.sz[dsu.find(a)]
}

// PairDist represents a pair with its squared distance
type PairDist struct {
	dist int
	i, j int
}

func readPoints(path string) ([]Point, error) {
	file, err := os.Open(path)
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
		if len(parts) != 3 {
			continue
		}
		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		z, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		points = append(points, Point{x, y, z})
	}
	return points, scanner.Err()
}

func squaredDistance(a, b Point) int {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

func solvePart1(path string, kPairs int) int {
	pts, err := readPoints(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return 0
	}

	n := len(pts)
	if n < 2 {
		fmt.Println("Not enough points.")
		return 0
	}

	var selected []PairDist

	// Build all pairs and sort
	totalPairs := n * (n - 1) / 2
	pairs := make([]PairDist, 0, totalPairs)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := squaredDistance(pts[i], pts[j])
			pairs = append(pairs, PairDist{d, i, j})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].dist < pairs[j].dist
	})
	if len(pairs) < kPairs {
		selected = pairs
	} else {
		selected = pairs[:kPairs]
	}

	// Union the selected pairs
	dsu := NewDSU(n)
	for _, pair := range selected {
		dsu.union(pair.i, pair.j)
	}

	// Count component sizes
	roots := make(map[int]int)
	for i := 0; i < n; i++ {
		root := dsu.find(i)
		roots[root]++
	}

	// Get top 3 sizes
	sizes := make([]int, 0, len(roots))
	for _, count := range roots {
		sizes = append(sizes, count)
	}
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	top3 := sizes
	if len(sizes) > 3 {
		top3 = sizes[:3]
	}

	// Calculate product
	prod := 1
	for _, s := range top3 {
		prod *= s
	}

	fmt.Printf("Number of points: %d\n", n)
	fmt.Printf("Total pairs considered (selected): %d\n", len(selected))
	fmt.Printf("Top 3 component sizes: %v\n", top3)
	fmt.Printf("Answer (product of top 3): %d\n", prod)

	return prod
}

func solvePart2(path string) int {
	pts, err := readPoints(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return 0
	}

	n := len(pts)
	if n < 2 {
		fmt.Println("Need at least two points")
		return 0
	}

	// Prim's algorithm to find MST and track largest edge
	inTree := make([]bool, n)
	bestDist := make([]int, n)
	parent := make([]int, n)

	// Initialize
	const largeNum = 1000000000
	for i := 0; i < n; i++ {
		bestDist[i] = largeNum
		parent[i] = -1
	}

	// Start from vertex 0
	bestDist[0] = 0

	maxEdgeW := -1
	maxEdgeU := -1
	maxEdgeV := -1

	// Build MST
	for iter := 0; iter < n; iter++ {
		// Find next vertex with minimum distance not yet in tree
		u := -1
		uDist := largeNum
		for i := 0; i < n; i++ {
			if !inTree[i] && bestDist[i] < uDist {
				uDist = bestDist[i]
				u = i
			}
		}

		if u == -1 {
			break
		}

		// Add u to tree
		inTree[u] = true

		// Track maximum edge
		if parent[u] != -1 {
			w := bestDist[u]
			if w > maxEdgeW {
				maxEdgeW = w
				maxEdgeU = parent[u]
				maxEdgeV = u
			}
		}

		// Update distances for remaining vertices
		for v := 0; v < n; v++ {
			if inTree[v] {
				continue
			}
			d := squaredDistance(pts[u], pts[v])
			if d < bestDist[v] {
				bestDist[v] = d
				parent[v] = u
			}
		}
	}

	if maxEdgeU == -1 || maxEdgeV == -1 {
		fmt.Println("No MST edge found")
		return 0
	}

	// Calculate product of X coordinates
	xProd := pts[maxEdgeU].x * pts[maxEdgeV].x

	fmt.Printf("Index endpoints: %d %d\n", maxEdgeU, maxEdgeV)
	fmt.Printf("Squared distance of that edge: %d\n", maxEdgeW)
	fmt.Printf("Product of X coordinates: %d\n", xProd)

	return xProd
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_8")

	fmt.Println("=== Day 8 - Part 1 ===")
	solvePart1(inputPath, 1000)

	fmt.Println("\n=== Day 8 - Part 2 ===")
	solvePart2(inputPath)
}

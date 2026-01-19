package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// parseLines parses input lines into an adjacency list
func parseLines(lines []string) map[string][]string {
	adj := make(map[string][]string)

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		if !strings.Contains(line, ":") {
			// Single token with no outputs
			node := line
			if _, exists := adj[node]; !exists {
				adj[node] = []string{}
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		node := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])

		// Parse outputs
		var outs []string
		if right != "" {
			for _, t := range strings.Fields(right) {
				if t != "" {
					outs = append(outs, t)
				}
			}
		}

		adj[node] = append(adj[node], outs...)

		// Ensure all neighbor nodes exist
		for _, n := range outs {
			if _, exists := adj[n]; !exists {
				adj[n] = []string{}
			}
		}
	}

	return adj
}

// countPaths counts distinct simple paths from start to target
func countPaths(adj map[string][]string, start, target string) int {
	memo := make(map[string]int)
	visiting := make(map[string]bool)

	var dfs func(string) int
	dfs = func(u string) int {
		if u == target {
			return 1
		}
		if val, exists := memo[u]; exists {
			return val
		}
		if visiting[u] {
			// Found a cycle - do not count paths that revisit nodes
			return 0
		}

		visiting[u] = true
		total := 0
		for _, v := range adj[u] {
			total += dfs(v)
		}
		delete(visiting, u)

		memo[u] = total
		return total
	}

	return dfs(start)
}

// Part 2 implementation
func buildIndegrees(edges map[string][]string, nodes map[string]bool) map[string]int {
	indeg := make(map[string]int)
	for u := range edges {
		for _, v := range edges[u] {
			indeg[v]++
		}
	}
	return indeg
}

func solvePart2(edges map[string][]string, nodes map[string]bool, start, end string) int {
	if !nodes[start] || !nodes[end] {
		return 0
	}

	indeg := buildIndegrees(edges, nodes)

	// DP with 4 masks: 0=none, 1=dac, 2=fft, 3=both
	dp := make(map[string][4]int)
	for n := range nodes {
		dp[n] = [4]int{0, 0, 0, 0}
	}

	// Set start mask
	startMask := 0
	if start == "dac" {
		startMask |= 1
	}
	if start == "fft" {
		startMask |= 2
	}

	startDP := dp[start]
	startDP[startMask] = 1
	dp[start] = startDP

	// Topological sort BFS
	queue := []string{start}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		// Propagate DP to all children
		for _, v := range edges[u] {
			// Determine mask addition for child
			addMask := 0
			if v == "dac" {
				addMask |= 1
			}
			if v == "fft" {
				addMask |= 2
			}

			vDP := dp[v]
			uDP := dp[u]
			for m := 0; m < 4; m++ {
				vDP[m|addMask] += uDP[m]
			}
			dp[v] = vDP

			indeg[v]--
			if indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	return dp[end][3] // Only paths that visited both dac and fft
}

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Construct path to input file
	inputPath := filepath.Join(cwd, "input_day_11")

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Input file not found: %v\n", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Part 1
	adj := parseLines(lines)
	result := countPaths(adj, "you", "out")
	fmt.Printf("Part 1: %d\n", result)

	// Part 2
	// Build nodes set and edges for Part 2
	edges := make(map[string][]string)
	nodes := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		src := strings.TrimSuffix(parts[0], ":")
		var right []string
		if len(parts) > 1 {
			right = parts[1:]
		}

		edges[src] = right
		nodes[src] = true
		for _, r := range right {
			nodes[r] = true
		}
	}

	// Ensure 'out' exists
	if _, exists := edges["out"]; !exists {
		edges["out"] = []string{}
	}
	nodes["out"] = true

	result2 := solvePart2(edges, nodes, "svr", "out")
	fmt.Printf("Part 2: %d\n", result2)
}

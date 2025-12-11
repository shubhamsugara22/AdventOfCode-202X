#!/usr/bin/env python3
import sys
from collections import defaultdict

def parse_lines(lines):
    """Parse lines like 'aaa: you hhh' -> adjacency dict."""
    adj = defaultdict(list)
    for raw in lines:
        line = raw.strip()
        if not line:
            continue
        if ":" not in line:
            # tolerate lines with single token meaning no outputs
            node = line
            adj[node] = adj[node]  # ensure present
            continue
        left, right = line.split(":", 1)
        node = left.strip()
        # outputs split by whitespace (may be none)
        outs = [t for t in right.strip().split() if t]
        adj[node].extend(outs)
        # ensure all neighbor nodes exist in dict (optional)
        for n in outs:
            if n not in adj:
                adj[n] = adj[n]
    return adj

def count_paths(adj, start="you", target="out"):
    """
    Count distinct simple paths from start to target.
    Uses DFS with memoization. If a node is on the current recursion stack
    and encountered again, it represents a cycle -> return 0 for that branch.
    """
    memo = {}        # node -> number of paths from node to target
    visiting = set() # recursion stack

    def dfs(u):
        if u == target:
            return 1
        if u in memo:
            return memo[u]
        if u in visiting:
            # found a cycle; do not count paths that revisit nodes
            return 0
        visiting.add(u)
        total = 0
        for v in adj.get(u, ()):
            total += dfs(v)
        visiting.remove(u)
        memo[u] = total
        return total

    return dfs(start)

def main():
    path = sys.argv[1] if len(sys.argv) > 1 else "input_day_11"
    try:
        with open(path, "r") as f:
            lines = f.readlines()
    except FileNotFoundError:
        print(f"Input file not found: {path}", file=sys.stderr)
        sys.exit(1)

    adj = parse_lines(lines)
    result = count_paths(adj, start="you", target="out")
    print(result)

if __name__ == "__main__":
    main()

## TODO part 2 logic
from collections import *

with open('input_day_11') as f:
    lines = f.read().splitlines()

# Build graph
edges = defaultdict(list)
nodes = set()

for line in lines:
    left, *right = line.split()
    src = left[:-1]
    edges[src] = right
    nodes.add(src)
    for r in right:
        nodes.add(r)

# Ensure 'out' exists
edges.setdefault("out", [])

# Build indegree map (fresh copy each time)
def build_indegrees():
    indeg = defaultdict(int)
    for u in edges:
        for v in edges[u]:
            indeg[v] += 1
    return indeg

def solve_part2(start="svr", end="out"):
    if start not in nodes or end not in nodes:
        return 0

    indeg = build_indegrees()

    # 3 masks: 0=no special, 1=dac, 2=fft, 3=both
    dp = {n: [0,0,0,0] for n in nodes}

    start_mask = 0
    if start == "dac": start_mask |= 1
    if start == "fft": start_mask |= 2

    dp[start][start_mask] = 1

    # queue for topo BFS
    q = deque()

    # start's indegree must be 0 to be processed first
    # but if it's not, we still add it manually
    q.append(start)

    while q:
        u = q.popleft()

        # propagate DP to all children
        for v in edges[u]:
            # determine mask addition for child
            add_mask = 0
            if v == "dac": add_mask |= 1
            if v == "fft": add_mask |= 2

            for m in range(4):
                dp[v][m | add_mask] += dp[u][m]

            indeg[v] -= 1
            if indeg[v] == 0:
                q.append(v)

    return dp[end][3]   # only paths that visited both


# ---------- RUN ----------
print("Part 2:", solve_part2("svr", "out"))
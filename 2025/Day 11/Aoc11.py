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

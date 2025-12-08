#!/usr/bin/env python3
import math
import sys
from itertools import combinations
from collections import Counter
from heapq import nsmallest
from typing import List, Tuple

class DSU:
    def __init__(self, n):
        self.p = list(range(n))
        self.r = [0]*n
        self.sz = [1]*n
    def find(self, a):
        p = self.p
        while p[a] != a:
            # path compression (iterative)
            p[a] = p[p[a]]
            a = p[a]
        return a
    def union(self, a, b):
        ra = self.find(a)
        rb = self.find(b)
        if ra == rb:
            return False
        if self.r[ra] < self.r[rb]:
            ra, rb = rb, ra
        self.p[rb] = ra
        self.sz[ra] += self.sz[rb]
        if self.r[ra] == self.r[rb]:
            self.r[ra] += 1
        return True
    def size(self, a):
        return self.sz[self.find(a)]

def read_points(path: str) -> List[Tuple[int,int,int]]:
    pts = []
    with open(path, "r") as f:
        for line in f:
            s = line.strip()
            if not s:
                continue
            x,y,z = map(int, s.split(","))
            pts.append((x,y,z))
    return pts

def squared_distance(a: Tuple[int,int,int], b: Tuple[int,int,int]) -> int:
    dx = a[0] - b[0]
    dy = a[1] - b[1]
    dz = a[2] - b[2]
    return dx*dx + dy*dy + dz*dz

def solve(path: str, k_pairs: int = 1000, use_heap_if_big: bool = True):
    pts = read_points(path)
    n = len(pts)
    if n < 2:
        print("Not enough points.")
        return 0

    total_pairs = n*(n-1)//2
    # If the total number of pairs is small-ish, build and sort all pairs.
    # Otherwise use heapq.nsmallest to avoid full sort if desired.
    # We'll use threshold: if total_pairs <= 5_000_000 (approx) it's fine to build full list in memory for typical inputs.
    threshold = 5_000_000

    if total_pairs <= threshold or not use_heap_if_big:
        pairs = []
        for i in range(n):
            xi, yi, zi = pts[i]
            for j in range(i+1, n):
                d = (xi-pts[j][0])**2 + (yi-pts[j][1])**2 + (zi-pts[j][2])**2
                pairs.append((d, i, j))
        # sort ascending by distance
        pairs.sort(key=lambda t: t[0])
        selected = pairs[:min(k_pairs, len(pairs))]
    else:
        # Use nsmallest to get k_pairs smallest pairs without storing all pairs
        # Generate pairs lazily
        def gen_pairs():
            for i in range(n):
                xi, yi, zi = pts[i]
                for j in range(i+1, n):
                    d = (xi-pts[j][0])**2 + (yi-pts[j][1])**2 + (zi-pts[j][2])**2
                    yield (d, i, j)
        selected = nsmallest(k_pairs, gen_pairs(), key=lambda t: t[0])

    # Now union the selected pairs in order
    dsu = DSU(n)
    # selected is already in ascending order; if we used nsmallest it might not be fully sorted,
    # so sort selected to ensure correct order processing.
    selected.sort(key=lambda t: t[0])
    for d, i, j in selected:
        dsu.union(i, j)

    # compute sizes
    roots = [dsu.find(i) for i in range(n)]
    counts = Counter(roots)
    sizes = sorted(counts.values(), reverse=True)

    # take three largest sizes; if less than 3 available, multiply what's there
    top3 = sizes[:3]
    # multiply
    prod = 1
    for s in top3:
        prod *= s

    print("Number of points:", n)
    print("Total pairs considered (selected):", len(selected))
    print("Top 3 component sizes:", top3)
    print("Answer (product of top 3):", prod)
    return prod

if __name__ == "__main__":
    # change filename if needed
    input_path = "input_day_8"
    solve(input_path)

#!/usr/bin/env python3
from typing import List, Tuple
import math

def read_points(path: str) -> List[Tuple[int,int,int]]:
    pts = []
    with open(path, "r") as f:
        for line in f:
            s = line.strip()
            if not s:
                continue
            x,y,z = map(int, s.split(","))
            pts.append((x,y,z))
    return pts

def sqdist(a: Tuple[int,int,int], b: Tuple[int,int,int]) -> int:
    dx = a[0] - b[0]
    dy = a[1] - b[1]
    dz = a[2] - b[2]
    return dx*dx + dy*dy + dz*dz

def prim_last_edge_xproduct(pts: List[Tuple[int,int,int]]):
    n = len(pts)
    if n < 2:
        raise ValueError("Need at least two points")

    in_tree = [False] * n
    best_dist = [10**30] * n   # large initial distance
    parent = [-1] * n

    # start from vertex 0
    best_dist[0] = 0
    parent[0] = -1

    # to store the maximum-weight MST edge (squared distance) and its endpoints
    max_edge_w = -1
    max_edge_u = -1
    max_edge_v = -1

    for _ in range(n):
        # pick next not-in-tree vertex with minimal best_dist
        u = -1
        u_dist = 10**30
        for i in range(n):
            if not in_tree[i] and best_dist[i] < u_dist:
                u_dist = best_dist[i]
                u = i
        if u == -1:
            break  # disconnected, shouldn't happen in complete graph

        # add u to the tree
        in_tree[u] = True

        # if u has a parent (i.e., not the initial root), consider its connecting edge
        if parent[u] != -1:
            w = best_dist[u]   # squared distance between parent[u] and u
            # track maximum MST edge weight and endpoints
            if w > max_edge_w:
                max_edge_w = w
                max_edge_u = parent[u]
                max_edge_v = u

        # update best_dist for remaining vertices using u as a new tree node
        for v in range(n):
            if in_tree[v]:
                continue
            d = sqdist(pts[u], pts[v])
            if d < best_dist[v]:
                best_dist[v] = d
                parent[v] = u

    # after MST built, max_edge_u and max_edge_v are the endpoints of the largest MST edge
    if max_edge_u == -1 or max_edge_v == -1:
        raise RuntimeError("No MST edge found (unexpected)")

    # multiply X coordinates
    x_prod = pts[max_edge_u][0] * pts[max_edge_v][0]
    return x_prod, (max_edge_u, max_edge_v, max_edge_w)

if __name__ == "__main__":
    input_path = "input_day_8"   # change if needed
    pts = read_points(input_path)
    prod, info = prim_last_edge_xproduct(pts)
    u, v, w = info[0], info[1], info[2]
    print("Index endpoints:", info[0], info[1])
    print("Squared distance of that edge:", info[2])
    print("Product of X coordinates:", prod)

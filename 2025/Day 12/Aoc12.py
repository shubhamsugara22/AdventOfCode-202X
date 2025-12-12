import sys
import itertools
import re

# ------------------------------------------
# Parse Input
# -----------------------------------------

def parse_input(filename):
    with open(filename) as f:
        lines = [line.rstrip("\n") for line in f]

    shapes = {}
    regions = []

    i = 0
    # Parse shapes section
    while i < len(lines):
        line = lines[i].strip()

        # Detect beginning of regions section
        if re.match(r"^\d+x\d+:", line):
            break

        # Shape header
        m = re.match(r"^(\d+):$", line)
        if m:
            idx = int(m.group(1))
            i += 1
            grid = []
            # Read shape lines until blank or region header
            while i < len(lines) and lines[i] and not re.match(r"^\d+x\d+:", lines[i]) and not re.match(r"^\d+:$", lines[i]):
                grid.append(lines[i])
                i += 1
            shapes[idx] = grid
            continue

        i += 1

    # Parse regions section
    while i < len(lines):
        line = lines[i].strip()
        if not line:
            i += 1
            continue

        m = re.match(r"^(\d+)x(\d+):\s*(.*)$", line)
        if m:
            W = int(m.group(1))
            H = int(m.group(2))
            counts = list(map(int, m.group(3).split()))
            regions.append((W, H, counts))
        
        i += 1

    return shapes, regions

# ------------------------------------------
# Shape transforms (rotations + flips)
# ------------------------------------------

def rotations_and_flips(shape):
    g = [list(row) for row in shape]

    def rotate(grid):
        R, C = len(grid), len(grid[0])
        return [[grid[R-1-r][c] for r in range(R)] for c in range(C)]

    def flip(grid):
        return [row[::-1] for row in grid]

    seen = set()
    out = []

    cur = g
    for _ in range(4):
        for f in [cur, flip(cur)]:
            key = tuple(tuple(r) for r in f)
            if key not in seen:
                seen.add(key)
                out.append(f)
        cur = rotate(cur)

    return out

# Convert shape to list of occupied coords
def shape_cells(grid):
    out = []
    for r in range(len(grid)):
        for c in range(len(grid[0])):
            if grid[r][c] == "#":
                out.append((r, c))
    return out

# ------------------------------------------
# Algorithm X / Dancing Links (Exact Cover)
# ------------------------------------------

class DLX:
    def __init__(self, subsets, universe):
        # subsets: list of lists
        # universe: list of column names
        self.U = list(universe)
        self.col_index = {c:i for i,c in enumerate(self.U)}

        self.cols = [0]*len(self.U)
        for s in subsets:
            for item in s:
                self.cols[self.col_index[item]] += 1

        self.subsets = subsets
        self.used = [False]*len(subsets)

    def solve(self):
        # If every column is covered → success
        if all(c == 0 for c in self.cols):
            return True

        # Choose column with minimal remaining options
        c = min((i for i,v in enumerate(self.cols) if v>0), key=lambda x:self.cols[x])

        # Try each subset covering column c
        for si, subset in enumerate(self.subsets):
            if self.used[si]: continue
            if self.U[c] not in subset: continue

            # Try using this subset
            removed = []
            self.used[si] = True
            for item in subset:
                ci = self.col_index[item]
                if self.cols[ci] > 0:
                    self.cols[ci] -= 1
                    removed.append(ci)

            if all(x == 0 for x in self.cols):
                return True

            if self.solve():
                return True

            # Backtrack
            for ci in removed:
                self.cols[ci] += 1
            self.used[si] = False

        return False


# ------------------------------------------
# Build placements for region
# ------------------------------------------

def can_fit_region(W, H, counts, all_shapes):
    needed = []
    for idx, c in enumerate(counts):
        needed.extend([(idx, k) for k in range(c)])

    # Universe columns: all piece identifiers + all grid cells
    universe = []
    for (idx, inst) in needed:
        universe.append(("piece", idx, inst))
    for r in range(H):
        for c in range(W):
            universe.append(("cell", r, c))

    subsets = []

    # Generate all placements
    for (idx, inst) in needed:
        variants = all_shapes[idx]
        for shape in variants:
            cells = shape_cells(shape)
            maxr = max(r for r,_ in cells)
            maxc = max(c for _,c in cells)

            for dr in range(H - maxr):
                for dc in range(W - maxc):
                    covered = []
                    ok = True
                    for (rr, cc) in cells:
                        R = rr + dr
                        C = cc + dc
                        if R<0 or C<0 or R>=H or C>=W:
                            ok = False
                            break
                        covered.append(("cell", R, C))

                    if not ok:
                        continue

                    # subset = piece instance + covered cells
                    subsets.append([("piece", idx, inst)] + covered)

    dlx = DLX(subsets, universe)
    return dlx.solve()

# ------------------------------------------
# Main Part 1 solver
# ------------------------------------------

def solve_day12_part1(filename):
    shapes, regions = parse_input(filename)

    # Precompute all variants (rotations + flips)
    all_shapes = {}
    for idx, g in shapes.items():
        all_shapes[idx] = rotations_and_flips(g)

    count = 0
    for W, H, counts in regions:
        if can_fit_region(W, H, counts, all_shapes):
            count += 1

    return count

# ------------------------------------------
# Run
# ------------------------------------------

print("Part 1:", solve_day12_part1("input_day_12"))

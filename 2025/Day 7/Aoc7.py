from typing import List, Set, Tuple

def count_splits(grid_lines: List[str]) -> int:
    # Normalize grid
    grid = [list(line.rstrip("\n")) for line in grid_lines]
    if not grid:
        return 0
    R = len(grid)
    C = max(len(row) for row in grid)
    # pad rows to equal width with spaces
    for i in range(R):
        if len(grid[i]) < C:
            grid[i].extend(" " * (C - len(grid[i])))

    # find source 'S'
    source = None
    for r in range(R):
        for c in range(C):
            if grid[r][c] == 'S':
                source = (r, c)
                break
        if source:
            break
    if source is None:
        raise ValueError("No source 'S' found in grid")

    # active beams as set of (row,col)
    active: Set[Tuple[int,int]] = {source}
    splits = 0
    seen_split_positions: Set[Tuple[int,int]] = set()

    while active:
        new_active: Set[Tuple[int,int]] = set()
        # process each beam: attempt to move one row down
        for (r, c) in active:
            nr = r + 1
            if nr >= R:
                # beam leaves the grid
                continue
            cell_below = grid[nr][c]
            if cell_below == '^':
                # Split occurs at (nr, c)
                if (nr, c) not in seen_split_positions:
                    splits += 1
                    seen_split_positions.add((nr, c))
                # spawn beams at immediate left and right of splitter (same row nr)
                left = (nr, c - 1)
                right = (nr, c + 1)
                if 0 <= left[1] < C:
                    # add left spawn position even if it's '^' or '.' or 'S' — beam sits there and will move down next step
                    new_active.add(left)
                if 0 <= right[1] < C:
                    new_active.add(right)
                # original beam does not continue downward beyond the '^'
            else:
                # move into the cell below (covers '.' and 'S' or any char not '^')
                new_active.add((nr, c))

        # deduplicate and continue
        active = new_active

    return splits


# Example usage:
if __name__ == "__main__":
    path = "input_day_7"   # change if needed
    with open(path, "r") as f:
        lines = [line.rstrip("\n") for line in f]
    result = count_splits(lines)
    print("Total splits:", result)

from typing import List, Set, Tuple

def count_timelines(grid_lines: List[str]) -> int:
    """
    Simulate the quantum tachyon manifold and return the total beam strength
    at the bottom row (number of beams that reach the exit).
    
    Key insight: Beams can merge and their strengths accumulate. When beams
    split at '^', the strength is divided and carried forward.
    """
    # Normalize grid
    grid = [list(line.rstrip("\n")) for line in grid_lines]
    if not grid:
        return 0
    R = len(grid)
    C = max(len(row) for row in grid) if grid else 0
    for i in range(R):
        if len(grid[i]) < C:
            grid[i].extend([" "] * (C - len(grid[i])))

    # Pad each row to width C
    for row in grid:
        while len(row) < C:
            row.append(" ")

    # Initialize beam strength matrix
    strength = [[0] * C for _ in range(R)]
    
    # Find source 'S' and initialize strength
    active_cols = set()  # columns with active beams
    for c in range(C):
        if grid[0][c] == "S":
            strength[0][c] = 1
            active_cols.add(c)

    # Process row by row
    for y in range(R - 1):
        next_active_cols = set()
        
        for x in active_cols:
            # Get cell at next row
            cell = grid[y + 1][x]
            
            # Check if cell above is a splitter
            cell_above = grid[y][x]
            is_below_splitter = (cell_above == "^")
            
            if cell == "^":
                # This row has a splitter: split the beam left and right
                left = x - 1
                right = x + 1
                
                if 0 <= left < C:
                    strength[y + 1][left] += strength[y][x]
                    next_active_cols.add(left)
                if 0 <= right < C:
                    strength[y + 1][right] += strength[y][x]
                    next_active_cols.add(right)
            else:
                # Regular cell or 'S': beam continues if not directly below a splitter
                if not is_below_splitter:
                    strength[y + 1][x] += strength[y][x]
                    next_active_cols.add(x)
        
        active_cols = next_active_cols

    # Return sum of strengths at the bottom row (beams reaching exit)
    return sum(strength[R - 1])


if __name__ == "__main__":
    # change path if your input file is elsewhere
    path = "input_day_7"
    with open(path, "r") as f:
        lines = [line.rstrip("\n") for line in f]
    result = count_timelines(lines)
    print("Number of timelines (Part 2):", result)

import re
from itertools import product

def parse_line(line):
    pattern = r"\[([.#]+)\](.*)\{.*\}"
    m = re.match(pattern, line.strip())
    lights = m.group(1)

    rest = m.group(2)
    buttons = []
    for part in re.findall(r"\((.*?)\)", rest):
        if part.strip() == "":
            buttons.append([])
        else:
            buttons.append(list(map(int, part.split(","))))
    return lights, buttons


def solve_machine(lights, buttons):
    n = len(lights)
    m = len(buttons)

    target = [1 if c == "#" else 0 for c in lights]

    # Build matrix B: size n x m, B[row][col]
    B = [[0]*m for _ in range(n)]
    for j, btn in enumerate(buttons):
        for bit in btn:
            B[bit][j] ^= 1

    # Gaussian elimination to find solution space
    row = 0
    where = [-1]*n

    for col in range(m):
        sel = -1
        for r in range(row, n):
            if B[r][col]:
                sel = r
                break
        if sel == -1:
            continue

        B[row], B[sel] = B[sel], B[row]
        target[row], target[sel] = target[sel], target[row]
        where[row] = col

        for r in range(n):
            if r != row and B[r][col]:
                for c in range(col, m):
                    B[r][c] ^= B[row][c]
                target[r] ^= target[row]

        row += 1

    # Check for inconsistency
    for r in range(row, n):
        if target[r]:
            return None  # no solution

    # Back-substitute general form
    x = [0]*m
    free_vars = []

    pivot_cols = set()
    for r, c in enumerate(where):
        if c != -1:
            pivot_cols.add(c)

    for col in range(m):
        if col not in pivot_cols:
            free_vars.append(col)

    best = None

    # Try all free variable combinations (usually small)
    for bits in product([0,1], repeat=len(free_vars)):
        x_try = x[:]

        for idx, col in enumerate(free_vars):
            x_try[col] = bits[idx]

        # Determine pivot variable values
        for r, col in enumerate(where):
            if col != -1:
                val = target[r]
                for c in range(col+1, m):
                    if B[r][c]:
                        val ^= x_try[c]
                x_try[col] = val

        weight = sum(x_try)
        if best is None or weight < best:
            best = weight

    return best


def solve_day10_part1(filename):
    total = 0
    with open(filename) as f:
        for line in f:
            if line.strip():
                lights, buttons = parse_line(line)
                presses = solve_machine(lights, buttons)
                total += presses
    return total


# Example run
if __name__ == "__main__":
    print("Day 10 Part 1:", solve_day10_part1("input_day_10"))

## Part 2
import re
import itertools
import numpy as np


def parse_line_part2(line):
    pattern = r"\[(.*?)\](.*)\{(.*?)\}"
    m = re.match(pattern, line.strip())
    button_section = m.group(2)
    target_str = m.group(3)

    buttons = []
    for part in re.findall(r"\((.*?)\)", button_section):
        part = part.strip()
        if not part:
            buttons.append([])
        else:
            buttons.append(list(map(int, part.split(","))))

    targets = list(map(int, target_str.split(",")))
    return buttons, targets


def solve_machine_part2(buttons, target):
    t = np.array(target, dtype=float)
    m = len(buttons)
    n = len(target)

    # Build matrix B: n counters × m buttons
    B = np.zeros((n, m), float)
    for j, btn in enumerate(buttons):
        for bit in btn:
            B[bit, j] += 1

    # Solve least-squares to get starting integer solution
    x0 = np.linalg.lstsq(B, t, rcond=None)[0]
    x0 = np.round(x0).astype(int)

    # Compute nullspace correctly
    u, s, vh = np.linalg.svd(B, full_matrices=True)
    tol = 1e-10

    # vh is (m, m)
    # singular values only exist for min(n, m)
    null_indices = []
    for i in range(vh.shape[0]):
        if i >= len(s) or s[i] < tol:
            null_indices.append(i)

    nullspace = vh[null_indices, :]

    # Convert each basis vector to integer-ish direction
    basis = []
    for v in nullspace:
        v = np.round(v).astype(int)
        if np.any(v != 0):
            basis.append(v)

    # If no nullspace, check direct solution
    if not basis:
        if np.all(x0 >= 0) and np.allclose(B @ x0, t):
            return int(np.sum(x0))
        else:
            return float("inf")

    # Search integer combinations
    best = None
    SEARCH = range(-20, 21)

    for ks in itertools.product(SEARCH, repeat=len(basis)):
        x = x0.copy()
        for k, b in zip(ks, basis):
            x += k * b

        if np.any(x < 0):
            continue
        if np.allclose(B @ x, t):
            total = int(np.sum(x))
            if best is None or total < best:
                best = total

    return best


def solve_day10_part2(filename):
    machines = []
    with open(filename) as f:
        lines = f.read().strip().split("\n")
    
    i = 0
    while i < len(lines):
        if lines[i].startswith("Button A"):
            ax, ay = map(int, lines[i].split(":")[1].split(","))
            bx, by = map(int, lines[i+1].split(":")[1].split(","))
            px, py = map(int, lines[i+2].split(":")[1].split(","))
            machines.append((ax, ay, bx, by, px + 10000000000000, py + 10000000000000))
            i += 3
        i += 1
    
    total = 0
    for (ax, ay, bx, by, px, py) in machines:
        det = ax * by - ay * bx
        if det == 0:
            continue
        
        A_num = px * by - py * bx
        B_num = ax * py - ay * px
        
        if A_num % det != 0 or B_num % det != 0:
            continue
        
        A = A_num // det
        B = B_num // det
        
        if A >= 0 and B >= 0:
            total += 3 * A + B
    
    return total



if __name__ == "__main__":
    print("Day 10 Part 2:", solve_day10_part2("input_day_10"))

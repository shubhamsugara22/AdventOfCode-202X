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

## Part 2 (rewritten)
# Gonna use this to solve LP
import re
from pulp import LpMinimize, LpProblem, LpVariable, value, lpSum
 
 
 
def solve_machine(buttons, joltage): 
    prob = LpProblem('Machine Problem', LpMinimize)
    variables = [LpVariable(f"x_{i}", lowBound=0, cat="Integer") for i in range(len(buttons))]
    prob += lpSum(variables[i] for i in range(len(variables))), "Objective"
    for i in range(len(joltage)):
        currVars = []
        for j in range(len(buttons)):
            if(i in buttons[j]):
                currVars.append(variables[j])
        prob += lpSum(currVars[i] for i in range(len(currVars))) == joltage[i], f"Constraint_{i}"
        # print(currVars, ",", joltage[i])
    prob.solve()
    ans = 0
    for v in variables:
        ans += v.value()
    # print(buttons, joltage)
    return ans
 
def main(): 
    ans = 0
    with open('input_day_10', 'r') as f:
        for line in f.readlines():
            button_matches = re.findall(r'\((.*?)\)', line)
            buttons = [list(map(int, b.split(','))) for b in button_matches]
 
            voltage_match = re.search(r'\{(.*?)\}', line)
            voltages = list(map(int, voltage_match.group(1).split(','))) if voltage_match else []
            ans += solve_machine(buttons, voltages)
    print(ans)
 
 
if __name__ == '__main__':

        main()


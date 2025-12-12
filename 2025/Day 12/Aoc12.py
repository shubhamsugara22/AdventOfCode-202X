def parse_input(filename):
    shapes = {}
    regions = []
    with open(filename) as f:
        lines = [l.rstrip() for l in f]

    i = 0
    # parse shapes
    while i < len(lines):
        if not lines[i]:
            i += 1
            continue
        if ":" not in lines[i]:
            break
        left = lines[i].split(":")[0]
        if not left.isdigit():
            break
        idx = int(left)
        i += 1
        grid = []
        while i < len(lines) and lines[i]:
            grid.append(lines[i])
            i += 1
        shapes[idx] = grid
        i += 1

    # parse regions
    while i < len(lines):
        if not lines[i]:
            i += 1
            continue
        prefix, nums = lines[i].split(":")
        W, H = map(int, prefix.split("x"))
        counts = list(map(int, nums.split()))
        regions.append(((W, H), counts))
        i += 1

    return shapes, regions


def solve_day12_part1(filename):
    shapes, regions = parse_input(filename)

    # compute area of each shape
    shape_area = {i: sum(row.count("#") for row in grid)
                  for i, grid in shapes.items()}

    good = 0
    for (W, H), counts in regions:
        required = sum(shape_area[i] * counts[i] for i in range(len(counts)))
        available = W * H
        if required <= available:
            good += 1

    return good


print("Part 1:", solve_day12_part1("input_day_12"))

def solve_day9(points):
    """
    points: list of (x, y) tuples representing red tile positions.
    returns: largest rectangle area using two red tiles as opposite corners.
    """
    n = len(points)
    max_area = 0

    # Compare all pairs of points
    for i in range(n):
        x1, y1 = points[i]
        for j in range(i + 1, n):
            x2, y2 = points[j]

            # Must form a valid rectangle (different x and y)
            if x1 != x2 and y1 != y2:
                width = abs(x1 - x2) + 1
                height = abs(y1 - y2) + 1
                area = width * height
                max_area = max(max_area, area)

    return max_area

def read_input(filename):
    points = []
    with open(filename) as f:
        for line in f:
            x, y = map(int, line.strip().split(","))
            points.append((x, y))
    return points

points = read_input("input_day_9")
print("Answer:", solve_day9(points))

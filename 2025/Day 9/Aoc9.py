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
print("Answer part 1 :", solve_day9(points))


#!/usr/bin/env python3
"""
Day 9 Part 2 — Robust, memory-safe solution.

Reads red tile coordinates from 'input_day_9' (one "x,y" per line).
Prints the largest rectangle area that can be formed with opposite corners
on red tiles, and whose interior (and edges) contain only red or green tiles,
where green is defined as the polygon interior formed by connecting red tiles
in order with horizontal/vertical segments (wrapping).
"""

from typing import List, Tuple
import sys

INPUT_PATH = "input_day_9"

# ---------- I/O ----------
def read_points(path: str) -> List[Tuple[int,int]]:
    pts = []
    with open(path, "r") as f:
        for line in f:
            s = line.strip()
            if not s:
                continue
            x,y = map(int, s.split(","))
            pts.append((x,y))
    return pts

# ---------- Geometry helpers ----------
def orient(a: Tuple[int,int], b: Tuple[int,int], c: Tuple[int,int]) -> int:
    return (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])

def on_segment(a: Tuple[int,int], b: Tuple[int,int], c: Tuple[int,int]) -> bool:
    return min(a[0],b[0]) <= c[0] <= max(a[0],b[0]) and min(a[1],b[1]) <= c[1] <= max(a[1],b[1])

def segments_proper_intersect(a: Tuple[int,int], b: Tuple[int,int],
                              c: Tuple[int,int], d: Tuple[int,int]) -> bool:
    """
    Return True only if segments ab and cd have a proper (non-collinear) intersection.
    Collinear touching/overlaps are considered NOT "proper" and thus allowed.
    """
    o1 = orient(a,b,c)
    o2 = orient(a,b,d)
    o3 = orient(c,d,a)
    o4 = orient(c,d,b)
    # proper intersection requires strict sign differences
    if (o1 > 0 and o2 < 0 or o1 < 0 and o2 > 0) and (o3 > 0 and o4 < 0 or o3 < 0 and o4 > 0):
        return True
    return False

def point_on_segment(px:int, py:int, ax:int, ay:int, bx:int, by:int) -> bool:
    # integer collinearity and within bounding box
    if (bx-ax)*(py-ay) != (by-ay)*(px-ax):
        return False
    return min(ax,bx) <= px <= max(ax,bx) and min(ay,by) <= py <= max(ay,by)

def point_in_polygon(x: int, y: int, poly: List[Tuple[int,int]]) -> bool:
    """
    Ray-casting: returns True if point is inside or exactly on polygon edge.
    For robustness on integer coordinates:
      - first check if point lies on any polygon edge (on-edge => True)
      - then do standard even-odd ray casting
    """
    n = len(poly)
    # On-edge test
    for i in range(n):
        x1,y1 = poly[i]
        x2,y2 = poly[(i+1)%n]
        if point_on_segment(x,y,x1,y1,x2,y2):
            return True
    inside = False
    for i in range(n):
        x1,y1 = poly[i]
        x2,y2 = poly[(i+1)%n]
        # check if edge crosses horizontal ray at y
        if (y1 > y) != (y2 > y):
            # compute x coordinate of intersection (float ok)
            xinters = x1 + (y - y1) * (x2 - x1) / (y2 - y1)
            if xinters >= x:
                inside = not inside
    return inside

# ---------- rectangle validity ----------
def rectangle_valid(p1: Tuple[int,int], p2: Tuple[int,int],
                    polygon: List[Tuple[int,int]], poly_bbox: Tuple[int,int,int,int]) -> bool:
    """
    Return True if the rectangle with diagonal corners p1 and p2
    lies entirely inside-or-on-boundary of polygon and does not
    have any proper crossing with the polygon boundary.
    """
    x1,y1 = p1; x2,y2 = p2
    if x1 == x2 or y1 == y2:
        return False
    minx, maxx = min(x1,x2), max(x1,x2)
    miny, maxy = min(y1,y2), max(y1,y2)
    bbox_minx, bbox_maxx, bbox_miny, bbox_maxy = poly_bbox
    # Quick bbox reject: rectangle must be within polygon bbox
    if minx < bbox_minx or maxx > bbox_maxx or miny < bbox_miny or maxy > bbox_maxy:
        return False

    corners = [(x1,y1),(x1,y2),(x2,y1),(x2,y2)]
    # All rectangle corners must be inside or on polygon boundary
    for cx,cy in corners:
        if not point_in_polygon(cx,cy,polygon):
            return False

    # Check that no rectangle edge properly crosses polygon edges.
    # Collinear overlaps and touching at vertices are allowed (these represent green/red tiles).
    rect_edges = [ (corners[0],corners[1]), (corners[1],corners[3]),
                   (corners[3],corners[2]), (corners[2],corners[0]) ]
    n = len(polygon)
    # For each rectangle edge, check against polygon edges for a proper intersection
    for a,b in rect_edges:
        # small bbox prune for this rect edge
        axmin, axmax = min(a[0],b[0]), max(a[0],b[0])
        aymin, aymax = min(a[1],b[1]), max(a[1],b[1])
        for i in range(n):
            c = polygon[i]; d = polygon[(i+1)%n]
            # Quick bbox overlap test: if bounding boxes don't overlap, skip
            if axmax < min(c[0],d[0]) or axmin > max(c[0],d[0]) or aymax < min(c[1],d[1]) or aymin > max(c[1],d[1]):
                continue
            if segments_proper_intersect(a,b,c,d):
                # proper crossing -> invalid rectangle
                return False
            # otherwise collinear/touching is OK, do not reject
    return True

# ---------- Main solve ----------
def solve_day9_part2_file(path: str) -> int:
    points = read_points(path)
    if not points:
        return 0
    polygon = points[:]  # polygons guaranteed by puzzle input
    xs = [p[0] for p in polygon]; ys = [p[1] for p in polygon]
    bbox = (min(xs), max(xs), min(ys), max(ys))
    n = len(points)
    best = 0
    best_pair = None

    # Heuristic ordering: try pairs that can produce large area first (helps prune)
    # Build index list sorted by x then y
    idxs = list(range(n))
    # You can sort by x or by y - choose x then y
    idxs.sort(key=lambda i: (points[i][0], points[i][1]))

    # Iterate pairs but prune by potential_area <= best
    for a_i in range(n):
        i = idxs[a_i]
        x1,y1 = points[i]
        for b_i in range(a_i+1, n):
            j = idxs[b_i]
            x2,y2 = points[j]
            if x1 == x2 or y1 == y2:
                continue
            potential_area = (abs(x1-x2)+1) * (abs(y1-y2)+1)
            if potential_area <= best:
                continue
            # quick bbox check and then exact geometry
            if rectangle_valid(points[i], points[j], polygon, bbox):
                best = potential_area
                best_pair = (points[i], points[j])
    # Done
    print("Best corners:", best_pair)
    return best

# ---------- Run ----------
if __name__ == "__main__":
    ans = solve_day9_part2_file("input_day_9")
    print("Part 2 largest rectangle area (red+green):", ans)

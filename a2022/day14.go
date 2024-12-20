package a2022

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

func pointFromString(s string) image.Point {
	parts := strings.Split(s, ",")
	return image.Pt(mustInt(parts[0]), mustInt(parts[1]))
}

func day14(in io.Reader) (int, int) {
	scanner := bufio.NewScanner(in)
	var rockLines [][]image.Point
	minX, maxX, maxY := 1000, 0, 0
	for scanner.Scan() {
		rockLines = append(
			rockLines,
			mapValue(strings.Split(scanner.Text(), " -> "), func(s string) image.Point {
				c := pointFromString(s)
				minX, maxX, maxY = minInt(minX, c.X), maxInt(maxX, c.X), maxInt(maxY, c.Y)
				return c
			}),
		)
	}

	maxY += 2
	grid := make([][]string, maxY)
	for i := range grid {
		row := make([]string, maxX+100)
		for j := range row {
			row[j] = aElseB(i == maxY, "#", " ")
		}
		grid[i] = row
	}

	for _, line := range rockLines {
		var from *image.Point
		for _, to := range line {
			if from != nil {
				dX, dY := sign(to.X-from.X), sign(to.Y-from.Y)
				// diagonal lines don't exist in the input, so one of these signs is
				// always zero. if that isn't true, we could maybe fix this by doing
				// conditional incrementing inside the loop
				for x, y := from.X, from.Y; x != to.X+dX || y != to.Y+dY; x, y = x+dX, y+dY {
					grid[y][x] = "#"
				}
			}
			temp := to
			from = &temp
		}
	}

	canDrop := func(c image.Point) (dX int, falls bool) {
		if c.Y+1 == maxY {
			return 0, false
		}
		for _, dX = range []int{0, -1, 1} {
			if grid[c.Y+1][c.X+dX] == " " {
				return dX, true
			}
		}
		return 0, false
	}

	sandUnits := 0
sandDrop:
	for {
		start := image.Pt(500, 0)
		for dX, falls := canDrop(start); falls; dX, falls = canDrop(start) {
			start.X, start.Y = start.X+dX, start.Y+1
			// pretend there's an abyss instead of the floor
			if start.Y == maxY-2 {
				break sandDrop
			}
		}
		grid[start.Y][start.X] = "o"
		sandUnits++
		if start.X == 500 && start.Y == 0 {
			break
		}
	}

	if shouldLog {
		for _, row := range grid {
			fmt.Printf("%s\n", strings.Join(row[minX-10:], ""))
		}
	}

	withFloor := sandUnits
	for {
		start := image.Pt(500, 0)
		for dX, falls := canDrop(start); falls; dX, falls = canDrop(start) {
			start.X, start.Y = start.X+dX, start.Y+1
		}
		grid[start.Y][start.X] = "o"
		withFloor++
		if start.X == 500 && start.Y == 0 {
			break
		}
	}

	if shouldLog {
		for _, row := range grid {
			fmt.Printf("%s\n", strings.Join(row[minX-10:], ""))
		}
	}

	return sandUnits, withFloor
}

package a2022

import (
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
)

type coord3d struct {
	X, Y, Z int
}

func day18(in io.Reader) (int, int) {
	var coords []coord3d
	if err := gocsv.UnmarshalCSVWithoutHeaders(csv.NewReader(in), &coords); err != nil {
		panic(err)
	}

	space := make([][][]bool, 20)
	for z := range space {
		col := make([][]bool, 20)
		for y := range col {
			col[y] = make([]bool, 20)
		}
		space[z] = col
	}

	contactSides := func(c coord3d) (sides int) {
		if c.Z < 19 {
			sides += aElseB(space[c.Z+1][c.Y][c.X], 1, 0)
		}
		if c.Z > 0 {
			sides += aElseB(space[c.Z-1][c.Y][c.X], 1, 0)
		}
		if c.Y < 19 {
			sides += aElseB(space[c.Z][c.Y+1][c.X], 1, 0)
		}
		if c.Y > 0 {
			sides += aElseB(space[c.Z][c.Y-1][c.X], 1, 0)
		}
		if c.X < 19 {
			sides += aElseB(space[c.Z][c.Y][c.X+1], 1, 0)
		}
		if c.X > 0 {
			sides += aElseB(space[c.Z][c.Y][c.X-1], 1, 0)
		}
		return sides
	}

	totalSurface := 0
	for _, c := range coords {
		newSurface := 6
		// for each side that touches, remove 1 from new droplet + 1 from touched spot
		newSurface -= 2 * contactSides(c)
		totalSurface += newSurface
		space[c.Z][c.Y][c.X] = true
	}

	anyTrue := func(cells []bool) bool {
		for _, c := range cells {
			if c {
				return true
			}
		}
		return false
	}
	exteriorSurface := totalSurface
	for z := range space {
		for y := range space[z] {
			for x := range space[z][y] {
				// we're only looking for water
				if space[z][y][x] {
					continue
				}
				// if it doesn't have contact with lava, it doesn't matter
				maybeSides := contactSides(coord3d{x, y, z})
				if maybeSides == 0 {
					continue
				}
				// if all 6 directions hit lava, we're trapped and need to exclude these sides
				beforeX, afterX := anyTrue(reverse(space[z][y][:x])), anyTrue(space[z][y][x+1:])
				col := mapValue(space[z], func(col []bool) bool { return col[x] })
				beforeY, afterY := anyTrue(reverse(col[:y])), anyTrue(col[y+1:])
				layer := mapValue(space, func(grid [][]bool) bool { return grid[y][x] })
				beforeZ, afterZ := anyTrue(reverse(layer[:z])), anyTrue(layer[z+1:])
				if beforeX && afterX && beforeY && afterY && beforeZ && afterZ {
					exteriorSurface -= maybeSides // we're on a water spot, so only the lava-side gets subtracted
				}
			}
		}
	}
	return totalSurface, exteriorSurface
}

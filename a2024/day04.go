package a2024

import (
	"bufio"
	"image"
	"io"
	"strings"
)

func day4(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	var lines [][]string
	for scanner.Scan() {
		lines = append(lines, strings.Split(scanner.Text(), ""))
	}

	finder := makeFinder(lines)

	var xmasCount int
	for y, row := range lines {
		for x, letter := range row {
			if letter == "X" {
				xmasCount += finder(image.Point{X: x, Y: y})
			}
		}
	}

	return xmasCount, 0
}

func makeFinder(lines [][]string) func(p image.Point) int {
	mX, mY := len(lines[0])-1, len(lines)-1
	sj := func(s []string) string {
		if len(s) != 4 {
			panic("isn't right length")
		}
		return strings.Join(s, "")
	}
	xmas := "XMAS"

	sliceFromPoint := func(p image.Point, dX, dY int) string {
		var ret []string
		for i := 0; i < 4; i++ {
			ret = append(ret, lines[p.Y][p.X])
			p.X += dX
			p.Y += dY
		}
		return sj(ret)
	}

	return func(p image.Point) int {
		found := 0
		if p.X+3 <= mX {
			found += aElseB(sliceFromPoint(p, 1, 0) == xmas, 1, 0)
			if p.Y >= 3 {
				found += aElseB(sliceFromPoint(p, 1, -1) == xmas, 1, 0)
			}
			if p.Y+3 <= mY {
				found += aElseB(sliceFromPoint(p, 1, 1) == xmas, 1, 0)
			}
		}
		if p.X >= 3 {
			found += aElseB(sliceFromPoint(p, -1, 0) == xmas, 1, 0)
			if p.Y >= 3 {
				found += aElseB(sliceFromPoint(p, -1, -1) == xmas, 1, 0)
			}
			if p.Y+3 <= mY {
				found += aElseB(sliceFromPoint(p, -1, 1) == xmas, 1, 0)
			}
		}
		if p.Y >= 3 {
			found += aElseB(sliceFromPoint(p, 0, -1) == xmas, 1, 0)
		}
		if p.Y+3 <= mY {
			found += aElseB(sliceFromPoint(p, 0, 1) == xmas, 1, 0)
		}

		return found
	}
}

package a2024

import (
	"bufio"
	"image"
	"io"
	"slices"
	"strings"
)

func day4(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	var lines [][]string
	for scanner.Scan() {
		lines = append(lines, strings.Split(scanner.Text(), ""))
	}

	finder := func(p image.Point) int {
		target := "XMAS"
		found := 0

		for dx := -1; dx < 2; dx++ {
			for dy := -1; dy < 2; dy++ {
				found += aElseB(sliceFromPoint(lines, len(target), p, dx, dy) == target, 1, 0)
			}
		}

		return found
	}

	crossFinder := func(p image.Point) int {
		targets := []string{"MAS", "SAM"}
		found := 0

		d1 := sliceFromPoint(lines, 3, image.Point{X: p.X - 1, Y: p.Y - 1}, 1, 1)
		d2 := sliceFromPoint(lines, 3, image.Point{X: p.X - 1, Y: p.Y + 1}, 1, -1)

		if slices.Contains(targets, d1) && slices.Contains(targets, d2) {
			found++
		}
		return found
	}

	xmasCount := 0
	crossCount := 0
	for y, row := range lines {
		for x, letter := range row {
			if letter == "X" {
				xmasCount += finder(image.Point{X: x, Y: y})
			}
			if letter == "A" {
				crossCount += crossFinder(image.Point{X: x, Y: y})
			}
		}
	}

	return xmasCount, crossCount
}

func sliceFromPoint(lines [][]string, length int, p image.Point, dX, dY int) string {
	mX, mY := len(lines[0])-1, len(lines)-1

	var ret []string
	for i := 0; i < length; i++ {
		if p.X < 0 || p.Y < 0 || p.X > mX || p.Y > mY {
			return strings.Join(ret, "")
		}
		ret = append(ret, lines[p.Y][p.X])
		p.X += dX
		p.Y += dY
	}

	return strings.Join(ret, "")
}

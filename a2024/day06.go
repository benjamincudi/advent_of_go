package a2024

import (
	"bufio"
	"image"
	"io"
	"strings"
)

func day6(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	obstacles := map[image.Point]struct{}{}
	var currentLoc image.Point
	maxY := -1
	maxX := 0
	for scanner.Scan() {
		maxY++
		row := strings.Split(scanner.Text(), "")
		maxX = len(row) - 1
		for currentX, char := range row {
			if char == "#" {
				obstacles[image.Point{X: currentX, Y: maxY}] = struct{}{}
			}
			if char == "^" {
				currentLoc = image.Point{X: currentX, Y: maxY}
			}
		}

	}

	visitedPoints := map[image.Point]struct{}{currentLoc: {}}
	walkingDirections := [][]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	walkingIndex := 0
	for {
		d := walkingDirections[walkingIndex]
		nextStep := image.Point{X: currentLoc.X + d[0], Y: currentLoc.Y + d[1]}
		if _, blocked := obstacles[nextStep]; blocked {
			walkingIndex = (walkingIndex + 1) % len(walkingDirections)
			continue
		}
		if nextStep.X < 0 || nextStep.X > maxX || nextStep.Y < 0 || nextStep.Y > maxY {
			break
		}
		visitedPoints[nextStep] = struct{}{}
		currentLoc = nextStep
	}

	return len(visitedPoints), 0
}

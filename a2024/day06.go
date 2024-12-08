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
	var startingLoc image.Point
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
				startingLoc = image.Point{X: currentX, Y: maxY}
			}
		}

	}

	type direction image.Point
	walkingDirections := []direction{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	getDir := func(index int) direction { return walkingDirections[index%4] }

	walkingIndex := 0
	currentLoc := startingLoc
	visitedPoints := map[image.Point]struct{}{startingLoc: {}}
	for {
		d := getDir(walkingIndex)
		nextStep := image.Point{X: currentLoc.X + d.X, Y: currentLoc.Y + d.Y}
		if _, blocked := obstacles[nextStep]; blocked {
			walkingIndex++
			continue
		}
		if nextStep.X < 0 || nextStep.X > maxX || nextStep.Y < 0 || nextStep.Y > maxY {
			break
		}
		currentLoc = nextStep
		visitedPoints[currentLoc] = struct{}{}
	}

	hasCycle := func() bool {
		currentLoc = startingLoc
		cycleWalkingIndex := 0
		cycleVisits := map[image.Point]map[direction]struct{}{startingLoc: {getDir(cycleWalkingIndex): {}}}
		for {
			d := getDir(cycleWalkingIndex)
			nextStep := image.Point{X: currentLoc.X + d.X, Y: currentLoc.Y + d.Y}
			if _, blocked := obstacles[nextStep]; blocked {
				cycleWalkingIndex++
				continue
			}
			if nextStep.X < 0 || nextStep.X > maxX || nextStep.Y < 0 || nextStep.Y > maxY {
				return false
			}
			if dirs, visited := cycleVisits[nextStep]; visited {
				if _, sameDir := dirs[d]; sameDir {
					return true
				}
			} else {
				cycleVisits[nextStep] = map[direction]struct{}{}
			}
			cycleVisits[nextStep][d] = struct{}{}
			currentLoc = nextStep
		}
	}

	cycleCount := 0
	for theoreticalObstruction := range visitedPoints {
		if theoreticalObstruction == startingLoc {
			continue
		}
		obstacles[theoreticalObstruction] = struct{}{}
		if hasCycle() {
			cycleCount++
		}
		delete(obstacles, theoreticalObstruction)
	}

	return len(visitedPoints), cycleCount
}

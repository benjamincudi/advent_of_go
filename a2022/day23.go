package a2022

import (
	"bufio"
	"image"
	"io"
	"strings"
)

func day23(in io.Reader) (int, int) {
	type elf struct {
		loc           image.Point
		nextDirection walkDirection
	}

	checkShouldMove := func(e *elf, grid map[image.Point]bool) bool {
		return grid[e.loc.Add(image.Pt(-1, 0))] ||
			grid[e.loc.Add(image.Pt(-1, 1))] ||
			grid[e.loc.Add(image.Pt(0, 1))] ||
			grid[e.loc.Add(image.Pt(1, 1))] ||
			grid[e.loc.Add(image.Pt(1, 0))] ||
			grid[e.loc.Add(image.Pt(1, -1))] ||
			grid[e.loc.Add(image.Pt(0, -1))] ||
			grid[e.loc.Add(image.Pt(-1, -1))]
	}

	checkStepDirection := func(e *elf, d walkDirection, grid map[image.Point]bool) bool {
		w := grid[e.loc.Add(image.Pt(-1, 0))]
		nw := grid[e.loc.Add(image.Pt(-1, 1))]
		n := grid[e.loc.Add(image.Pt(0, 1))]
		ne := grid[e.loc.Add(image.Pt(1, 1))]
		east := grid[e.loc.Add(image.Pt(1, 0))]
		se := grid[e.loc.Add(image.Pt(1, -1))]
		s := grid[e.loc.Add(image.Pt(0, -1))]
		sw := grid[e.loc.Add(image.Pt(-1, -1))]
		switch d {
		case left:
			return !w && !nw && !sw
		case right:
			return !east && !ne && !se
		case up:
			return !n && !nw && !ne
		case down:
			return !s && !sw && !se
		default:
			panic("unknown walking direction")
		}
	}

	scanner := bufio.NewScanner(in)
	currentRow := 0
	var elves []*elf
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "")
		for x, s := range row {
			if s == "#" {
				elves = append(elves, &elf{image.Pt(x, currentRow), up})
			}
		}
		// in order for the first row parsed to be the "top", decrease Y
		currentRow--
	}

	elfDirOrder := map[walkDirection][]walkDirection{
		up:    {up, down, left, right},
		down:  {down, left, right, up},
		left:  {left, right, up, down},
		right: {right, up, down, left},
	}
	dirToVector := map[walkDirection]image.Point{
		up: {0, 1}, left: {-1, 0},
		down: {0, -1}, right: {1, 0},
	}

	emptySpacesAt10, finalRound := 0, 0
	for round := 1; true; round++ {
		currentGrid := map[image.Point]bool{}
		for _, e := range elves {
			currentGrid[e.loc] = true
		}
		stationaryElves := 0
		proposedStep := map[image.Point][]*elf{}
		for _, e := range elves {
			if !checkShouldMove(e, currentGrid) {
				stationaryElves++
				continue
			}
			for _, d := range elfDirOrder[e.nextDirection] {
				if checkStepDirection(e, d, currentGrid) {
					dest := e.loc.Add(dirToVector[d])
					proposedStep[dest] = append(proposedStep[dest], e)
					break
				}
			}
		}
		if stationaryElves == len(elves) {
			finalRound = round
			break
		}
		for dest, group := range proposedStep {
			if len(group) == 1 {
				group[0].loc = dest
			}
		}
		for _, e := range elves {
			e.nextDirection = elfDirOrder[e.nextDirection][1]
		}
		if round == 10 {
			minX := minInt(mapValue(elves, func(e *elf) int { return e.loc.X })...)
			maxX := maxInt(mapValue(elves, func(e *elf) int { return e.loc.X })...)
			minY := minInt(mapValue(elves, func(e *elf) int { return e.loc.Y })...)
			maxY := maxInt(mapValue(elves, func(e *elf) int { return e.loc.Y })...)
			cg := map[image.Point]bool{}
			for _, e := range elves {
				cg[e.loc] = true
			}
			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {
					if !cg[image.Pt(x, y)] {
						emptySpacesAt10++
					}
				}
			}
		}
	}

	return emptySpacesAt10, finalRound
}

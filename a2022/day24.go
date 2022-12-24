package a2022

import (
	"bufio"
	"image"
	"io"
	"math"
	"strings"
)

func day24(in io.Reader) (int, int) {
	type blizzard struct {
		loc image.Point
		d   walkDirection
	}
	blizzardDirectionLookup := map[string]walkDirection{
		"<": left, ">": right,
		"^": up, "v": down,
	}
	dirToVector := map[walkDirection]image.Point{
		left: {-1, 0}, right: {1, 0},
		up: {0, 1}, down: {0, -1},
	}

	var storms []blizzard
	scanner := bufio.NewScanner(in)
	var rawRows []string
	for scanner.Scan() {
		// For the sake of pretty-printing, we'll build the grid in another loop
		// so that we can do a simpler `for _, row := range grid { fmt.Printf }` loop
		rawRows = append(rawRows, scanner.Text())
	}
	grid := make([][]int, len(rawRows))
	// These are _always_ in these positions: opposite corners, inset from the L/R edge but on the outer wall
	entrance, exit := image.Pt(1, len(rawRows)-1), image.Pt(len(strings.Split(rawRows[0], ""))-2, 0)

	distanceToPoint := func(from, to image.Point) int {
		d := from.Sub(to)
		return abs(d.X) + abs(d.Y)
	}

	for i, text := range rawRows {
		cells := strings.Split(text, "")
		row := make([]int, len(cells))
		for x, c := range cells {
			switch c {
			case "#":
				row[x] = -1
				continue
			case ".":
				// just open space, nothing special
			default:
				storms = append(storms, blizzard{image.Pt(x, len(rawRows)-(i+1)), blizzardDirectionLookup[c]})
			}
			// ignore blizzards for initial weight - every square COULD be open if it isn't a wall
			row[x] = 0
		}
		grid[len(rawRows)-1-i] = row // assign in reverse for pretty-printing and coordinates not being backwards
	}

	gridRect := image.Rect(0, 0, len(grid[0]), len(grid))
	blizzardRect := gridRect.Inset(1)
	getStormLocationAtMinute := func(b blizzard, minute int) image.Point {
		return b.loc.Add(dirToVector[b.d].Mul(minute)).Mod(blizzardRect)
	}

	gridAtMinute := func(minute int) [][]int {
		copied := make([][]int, len(grid))
		// copy each slice, otherwise each row is a reference still
		for i := range copied {
			copied[i] = make([]int, len(grid[i]))
			copy(copied[i], grid[i])
		}

		for _, b := range storms {
			loc := getStormLocationAtMinute(b, minute)
			// spots with blizzards are equivalent to walls, we can't go there
			// it doesn't matter how many storms are on the square
			copied[loc.Y][loc.X] = -1
		}

		return copied
	}

	getPossibleMoves := func(from image.Point, gridState [][]int) []image.Point {
		var choices []image.Point
		for _, next := range []image.Point{
			from.Add(dirToVector[up]), from.Add(dirToVector[down]),
			from.Add(dirToVector[left]), from.Add(dirToVector[right]),
			from,
		} {
			// Make sure we're not trying to step out of bounds
			// This only applies at the entrance+exit points
			if !next.In(gridRect) {
				continue
			}
			// blizzards and walls all have score -1, search handles finding the fastest route
			if gridState[next.Y][next.X] == 0 {
				choices = append(choices, next)
			}
		}

		return choices
	}

	makeTrip := func(from, to image.Point, startingMinute int) int {
		type queue struct {
			minute int
			loc    image.Point
		}
		var state queue

		// The cache is only relevant for the trip we're making, since we
		// won't time-travel backwards, and it is cheap to re-build a minute
		// that we may have discarded
		minuteGridCache := map[int]*[][]int{}
		searchQueue := []queue{{startingMinute, from}}
		seenStates := map[queue]bool{}
		best := math.MaxInt
		for len(searchQueue) > 0 {
			state, searchQueue = searchQueue[0], searchQueue[1:]
			if _, seen := seenStates[state]; seen {
				// we've already searched from here before
				continue
			}
			seenStates[state] = true
			if state.minute+distanceToPoint(state.loc, to) > best {
				// If we could ignore blizzards and walk straight there and still
				// would be slower, give up on this path
				continue
			}
			cache, ok := minuteGridCache[state.minute+1]
			if !ok {
				gridState := gridAtMinute(state.minute + 1)
				minuteGridCache[state.minute+1] = &gridState
				cache = &gridState
			}
			choices := getPossibleMoves(state.loc, *cache)
			for _, c := range choices {
				if c == to {
					best = minInt(best, state.minute+1)
				} else {
					searchQueue = append(searchQueue, queue{state.minute + 1, c})
				}
			}
		}
		return best
	}

	firstTrip := makeTrip(entrance, exit, 0)
	thirdTrip := makeTrip(entrance, exit, makeTrip(exit, entrance, firstTrip))

	return firstTrip, thirdTrip
}

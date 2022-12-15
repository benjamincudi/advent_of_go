package a2022

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type sensorBeaconPairing struct {
	S, B coordinates
}

type coveredRange struct{ min, max int }

func (left *coveredRange) equals(right coveredRange) bool {
	return left.min == right.min && left.max == right.max
}

func (left *coveredRange) intersects(right coveredRange) bool {
	// fully covered - left inside of right
	if (left.min >= right.min && left.max <= right.max) ||
		// fully covered - right inside of left
		(left.min <= right.min && left.max >= right.max) ||
		// min-overlapping
		(left.min <= right.min && left.max >= right.min) ||
		// max-overlapping
		(left.min <= right.max && left.max >= right.max) {
		return true
	}
	return false
}
func (left *coveredRange) combine(right coveredRange) {
	left.min, left.max = minInt(left.min, right.min), maxInt(left.max, right.max)
}

type rowCoverage struct{ disjointRanges []coveredRange }

func (rc *rowCoverage) addRange(cr coveredRange) {
	anyIntersect := false
	for i, left := range rc.disjointRanges {
		if left.intersects(cr) {
			rc.disjointRanges[i].combine(cr)
			anyIntersect = true
			break
		}
	}
	if !anyIntersect {
		rc.disjointRanges = append(rc.disjointRanges, cr)
	} else {
		rc.consolidateRanges()
	}
}

func (rc *rowCoverage) consolidateRanges() {
	if len(rc.disjointRanges) == 1 {
		return
	}
	var consolidated []coveredRange
	// sorting ensures we only have to walk through once
	sort.SliceStable(rc.disjointRanges, func(i, j int) bool {
		return rc.disjointRanges[i].min < rc.disjointRanges[j].min
	})
	for i := 0; i < len(rc.disjointRanges); {
		left := rc.disjointRanges[i]
		nextI := i + 1
		// keep expanding the left until there is no overlap
		for ; nextI < len(rc.disjointRanges) && left.intersects(rc.disjointRanges[nextI]); nextI++ {
			left.combine(rc.disjointRanges[nextI])
		}
		// skip ahead to the first set we didn't combine with, everything in-between is part of left now
		consolidated, i = append(consolidated, left), nextI
	}
	rc.disjointRanges = consolidated
}

func day15(in io.Reader, targetRow, upperBounds int) (int, int) {
	scanner := bufio.NewScanner(in)
	var pairs []sensorBeaconPairing
	for scanner.Scan() {
		tokens := filter(
			strings.Split(scanner.Text(), " "),
			func(t string) bool {
				return strings.HasPrefix(t, "x=") || strings.HasPrefix(t, "y=")
			},
		)
		vals := mapValue(tokens, func(s string) int {
			return mustInt(strings.TrimRight(strings.Split(s, "=")[1], ":,"))
		})
		pairs = append(pairs, sensorBeaconPairing{coordinates{vals[0], vals[1]}, coordinates{vals[2], vals[3]}})
	}

	occupiedX := map[int]bool{}
	for _, p := range pairs {
		dTotal := abs(p.S.Y-p.B.Y) + abs(p.S.X-p.B.X)
		rowDiff := abs(targetRow - p.S.Y)
		if dTotal <= rowDiff {
			continue
		}
		xRange := dTotal - rowDiff
		for x := p.S.X - xRange; x <= p.S.X+xRange; x++ {
			occupiedX[x] = true
		}
	}
	for _, p := range pairs {
		if p.B.Y == targetRow {
			delete(occupiedX, p.B.X)
		}
	}
	var tuningFrequency int
	for y := 0; y <= upperBounds; y++ {
		rc := rowCoverage{}
		for _, p := range pairs {
			dTotal := abs(p.S.Y-p.B.Y) + abs(p.S.X-p.B.X)
			rowDiff := abs(y - p.S.Y)
			if dTotal <= rowDiff {
				continue
			}
			xRange := dTotal - rowDiff
			starting := maxInt(p.S.X-xRange, 0)
			end := minInt(p.S.X+xRange, upperBounds)
			if starting == 0 && end == upperBounds {
				// fully covers the row by itself - start next row
				rc.disjointRanges = []coveredRange{fullyCovered}
				break
			}
			rc.addRange(coveredRange{starting, end})
			if len(rc.disjointRanges) == 1 && rc.disjointRanges[0].equals(coveredRange{0, upperBounds}) {
				// row is covered, stop checking it
				rc.disjointRanges = []coveredRange{fullyCovered}
				break
			}
		}
		// AoC guarantees there is only one space to find, the search is over when there are 2
		// Separately, the control case has more than one spot, so we need to exit after the first
		if len(rc.disjointRanges) == 2 {
			left := maxInt(rc.disjointRanges[0].min, rc.disjointRanges[1].min)
			missingX := left - 1
			if shouldLog {
				fmt.Printf("distress coords %d,%d\n", missingX, y)
			}
			tuningFrequency = (missingX * 4000000) + y
			break
		}
	}
	return len(occupiedX), tuningFrequency
}

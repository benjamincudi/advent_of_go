package a2024

import (
	"image"
	"io"
	"strings"
)

type antennaPoint struct {
	antennaLabel string
	antinodes    map[string]struct{}
}

func (p *antennaPoint) antinodeFor(label string) {
	if p.antinodes == nil {
		p.antinodes = map[string]struct{}{}
	}
	p.antinodes[label] = struct{}{}
}

func day8(r io.Reader) (int, int) {
	gridInfo := map[image.Point]*antennaPoint{}
	yMax := 0
	xMax := 0
	labelLocations := map[string][]image.Point{}
	mapFileLines(r, func(text string) {
		cells := strings.Split(text, "")
		xMax = len(cells)
		for x, label := range cells {
			p := image.Point{X: x, Y: yMax}
			if label != "." {
				gridInfo[p] = &antennaPoint{antennaLabel: label}
				if locs, ok := labelLocations[label]; !ok {
					labelLocations[label] = []image.Point{p}
				} else {
					labelLocations[label] = append(locs, p)
				}
			}
		}
		yMax++
	})

	gridBox := image.Rectangle{image.Point{0, 0}, image.Point{xMax, yMax}}

	for label, locs := range labelLocations {
		for i, left := range locs {
			for j := i + 1; j < len(locs); j++ {
				right := locs[j]
				distanceToJ := right.Sub(left)
				pastRight := right.Add(distanceToJ)
				distanceFromJ := left.Sub(right)
				beforeLeft := left.Add(distanceFromJ)
				if pastRight.In(gridBox) {
					if _, seen := gridInfo[pastRight]; !seen {
						n := &antennaPoint{}
						n.antinodeFor(label)
						gridInfo[pastRight] = n
					} else {
						gridInfo[pastRight].antinodeFor(label)
					}
				}
				if beforeLeft.In(gridBox) {
					if _, seen := gridInfo[beforeLeft]; !seen {
						n := &antennaPoint{}
						n.antinodeFor(label)
						gridInfo[beforeLeft] = n
					} else {
						gridInfo[beforeLeft].antinodeFor(label)
					}
				}
			}
		}
	}

	hasAntinode := 0
	for _, p := range gridInfo {
		if len(p.antinodes) > 0 {
			hasAntinode++
		}
	}

	for label, locs := range labelLocations {
		for i, left := range locs {
			for j := i + 1; j < len(locs); j++ {
				right := locs[j]
				gridInfo[left].antinodeFor(label)
				gridInfo[right].antinodeFor(label)
				distanceToJ := right.Sub(left)
				for pastRight := right.Add(distanceToJ); pastRight.In(gridBox); pastRight = pastRight.Add(distanceToJ) {
					if _, seen := gridInfo[pastRight]; !seen {
						n := &antennaPoint{}
						n.antinodeFor(label)
						gridInfo[pastRight] = n
					} else {
						gridInfo[pastRight].antinodeFor(label)
					}
				}
				distanceFromJ := left.Sub(right)
				for beforeLeft := left.Add(distanceFromJ); beforeLeft.In(gridBox); beforeLeft = beforeLeft.Add(distanceFromJ) {
					if _, seen := gridInfo[beforeLeft]; !seen {
						n := &antennaPoint{}
						n.antinodeFor(label)
						gridInfo[beforeLeft] = n
					} else {
						gridInfo[beforeLeft].antinodeFor(label)
					}
				}
			}
		}
	}

	resonanceAntinode := 0
	for _, p := range gridInfo {
		if len(p.antinodes) > 0 {
			resonanceAntinode++
		}
	}

	return hasAntinode, resonanceAntinode
}

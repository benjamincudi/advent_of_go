package a2022

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

type tetrisShape interface {
	getHeight() int
	maybeWind(wind int, g tetrisGrid)
	maybeFall(g tetrisGrid) bool
	stopFalling() []image.Point
}

type tetrisGrid struct {
	state [][]bool
}

func (tg *tetrisGrid) getHeight() int {
	return len(tg.state)
}

func (tg *tetrisGrid) addRows(height int) {
	for i := 0; i < height; i++ {
		tg.state = append(tg.state, make([]bool, 7))
	}
}

func (tg *tetrisGrid) isOpen(p image.Point) bool {
	if p.X < 0 || p.X > 6 || p.Y < 0 {
		return false
	}
	if p.Y >= len(tg.state) {
		return true
	}
	return !tg.state[p.Y][p.X]
}

func (tg *tetrisGrid) recordBlocks(s tetrisShape) {
	points := s.stopFalling()
	reqHeight := maxInt(mapValue(points, func(p image.Point) int { return p.Y })...)
	if reqHeight+1 > tg.getHeight() {
		tg.addRows(reqHeight - tg.getHeight() + 1)
	}
	for _, p := range points {
		if tg.state[p.Y][p.X] {
			panic("position already occupied")
		} else {
			tg.state[p.Y][p.X] = true
		}
	}
}

func (tg *tetrisGrid) String() string {
	var grid string
	for j := tg.getHeight() - 1; j >= 0; j-- {
		grid += fmt.Sprintf("%v\n", mapValue(tg.state[j], func(s bool) string { return aElseB(s, "#", ".") }))
	}
	return grid
}

func day17(in io.Reader) (int, int) {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	wind := mapValue(strings.Split(scanner.Text(), ""), func(s string) int { return aElseB(s == "<", -1, 1) })
	windLen := len(wind)
	tg := tetrisGrid{[][]bool{}}
	turn := 0
	type cacheEntry struct{ numRocks, currentHeight []int }
	cache := map[int]cacheEntry{}
	rocksToDrop := 1000000000000
	part1height, elidedHeight := 0, 0
	for i := 0; i < rocksToDrop; i++ {
		if i == 2022 {
			part1height = tg.getHeight()
		}

		var shape tetrisShape
		switch i % 5 {
		case 0:
			shape = makeTetrisLine(tg.getHeight())
			// lazy to ensure we don't find the cycle too soon and miss having an answer for part 1
			if i > 2022 && elidedHeight == 0 {
				if entry, cached := cache[turn%windLen]; cached {
					if len(entry.currentHeight) > 2 {
						histLen := len(entry.currentHeight) - 1
						if entry.currentHeight[histLen]-entry.currentHeight[histLen-1] == entry.currentHeight[histLen-1]-entry.currentHeight[histLen-2] {
							heightDiff := entry.currentHeight[histLen] - entry.currentHeight[histLen-2]
							cycleLength := entry.numRocks[histLen] - entry.numRocks[histLen-2]
							additionalCompleteCycles := (rocksToDrop-entry.numRocks[histLen-2])/cycleLength - 2
							elidedHeight = heightDiff * additionalCompleteCycles
							i += additionalCompleteCycles * cycleLength
							if i == rocksToDrop {
								return part1height, tg.getHeight() + elidedHeight
							}
						}
					}
					cache[turn%windLen] = cacheEntry{
						append(entry.numRocks, i),
						append(entry.currentHeight, tg.getHeight()),
					}
				} else {
					cache[turn%windLen] = cacheEntry{[]int{i}, []int{tg.getHeight()}}
				}
			}
		case 1:
			shape = makeTetrisPlus(tg.getHeight())
		case 2:
			shape = makeTetrisCorner(tg.getHeight())
		case 3:
			shape = makeTetrisCol(tg.getHeight())
		default:
			shape = makeTetrisSquare(tg.getHeight())
		}

		for {
			shape.maybeWind(wind[turn%windLen], tg)
			turn++
			if !shape.maybeFall(tg) {
				tg.recordBlocks(shape)
				break
			}
		}
		//fmt.Printf("%s\n", tg.String())
	}

	//fmt.Println("final grid\n\n ")
	//fmt.Printf("%s\n", tg.String())
	return part1height, tg.getHeight() + elidedHeight
}

type tetrisLine struct {
	leftBlock image.Point
}

func makeTetrisLine(gridHeight int) tetrisShape {
	return &tetrisLine{image.Pt(2, gridHeight+3)}
}

func (tl *tetrisLine) getHeight() int { return 1 }
func (tl *tetrisLine) maybeWind(wind int, tg tetrisGrid) {
	if tg.isOpen(aElseB(wind == -1, image.Pt(tl.leftBlock.X-1, tl.leftBlock.Y), image.Pt(tl.leftBlock.X+4, tl.leftBlock.Y))) {
		tl.leftBlock.X += wind
	}
}
func (tl *tetrisLine) maybeFall(tg tetrisGrid) bool {
	for x := 0; x < 4; x++ {
		if !tg.isOpen(image.Pt(tl.leftBlock.X+x, tl.leftBlock.Y-1)) {
			return false
		}
	}
	tl.leftBlock.Y -= 1
	return true
}
func (tl *tetrisLine) stopFalling() []image.Point {
	points := make([]image.Point, 4)
	for x := 0; x < 4; x++ {
		points[x] = image.Pt(tl.leftBlock.X+x, tl.leftBlock.Y)
	}
	return points
}

type tetrisCol struct {
	bottomBlock image.Point
}

func makeTetrisCol(gridHeight int) tetrisShape {
	return &tetrisCol{image.Pt(2, gridHeight+3)}
}

func (tl *tetrisCol) getHeight() int { return 4 }
func (tl *tetrisCol) maybeWind(wind int, tg tetrisGrid) {
	for y := 0; y < 4; y++ {
		if !tg.isOpen(image.Pt(tl.bottomBlock.X+wind, tl.bottomBlock.Y+y)) {
			return
		}
	}
	tl.bottomBlock.X += wind
}
func (tl *tetrisCol) maybeFall(tg tetrisGrid) bool {
	if tl.bottomBlock.Y == 0 {
		return false
	}
	if !tg.isOpen(image.Pt(tl.bottomBlock.X, tl.bottomBlock.Y-1)) {
		return false
	}
	tl.bottomBlock.Y -= 1
	return true
}
func (tl *tetrisCol) stopFalling() []image.Point {
	points := make([]image.Point, 4)
	for y := 0; y < 4; y++ {
		points[y] = image.Pt(tl.bottomBlock.X, tl.bottomBlock.Y+y)
	}
	return points
}

type tetrisSquare struct {
	bottomLeftBlock image.Point
}

func makeTetrisSquare(gridHeight int) tetrisShape {
	return &tetrisSquare{image.Pt(2, gridHeight+3)}
}

func (tl *tetrisSquare) getHeight() int { return 2 }
func (tl *tetrisSquare) maybeWind(wind int, tg tetrisGrid) {
	edge := aElseB(sign(wind) == -1, 0, 1)
	for y := 0; y < 2; y++ {
		if !tg.isOpen(image.Pt(tl.bottomLeftBlock.X+edge+wind, tl.bottomLeftBlock.Y+y)) {
			return
		}
	}
	tl.bottomLeftBlock.X += wind
}
func (tl *tetrisSquare) maybeFall(tg tetrisGrid) bool {
	if tl.bottomLeftBlock.Y == 0 {
		return false
	}
	for x := 0; x < 2; x++ {
		if !tg.isOpen(image.Pt(tl.bottomLeftBlock.X+x, tl.bottomLeftBlock.Y-1)) {
			return false
		}
	}
	tl.bottomLeftBlock.Y -= 1
	return true
}
func (tl *tetrisSquare) stopFalling() []image.Point {
	points := make([]image.Point, 4)
	for y := 0; y < 2; y++ {
		points[2*y] = image.Pt(tl.bottomLeftBlock.X, tl.bottomLeftBlock.Y+y)
		points[2*y+1] = image.Pt(tl.bottomLeftBlock.X+1, tl.bottomLeftBlock.Y+y)
	}
	return points
}

// The backwards L whose arm is too long
type tetrisCorner struct {
	bottomLeftBlock image.Point
}

func makeTetrisCorner(gridHeight int) tetrisShape {
	return &tetrisCorner{image.Pt(2, gridHeight+3)}
}

func (tl *tetrisCorner) getHeight() int { return 3 }
func (tl *tetrisCorner) maybeWind(wind int, tg tetrisGrid) {
	for y := 0; y < 3; y++ {
		dest := image.Pt(tl.bottomLeftBlock.X+aElseB(y == 0 && wind == -1, 0, 2)+wind, tl.bottomLeftBlock.Y+y)
		if !tg.isOpen(dest) {
			return
		}
	}
	tl.bottomLeftBlock.X += wind
}
func (tl *tetrisCorner) maybeFall(tg tetrisGrid) bool {
	for x := 0; x < 3; x++ {
		dest := image.Pt(tl.bottomLeftBlock.X+x, tl.bottomLeftBlock.Y-1)
		if !tg.isOpen(dest) {
			return false
		}
	}
	tl.bottomLeftBlock.Y -= 1
	return true
}
func (tl *tetrisCorner) stopFalling() []image.Point {
	points := make([]image.Point, 5)
	for x := 0; x < 3; x++ {
		points[x] = image.Pt(tl.bottomLeftBlock.X+x, tl.bottomLeftBlock.Y)
	}
	for y := 1; y < 3; y++ {
		points[y+2] = image.Pt(tl.bottomLeftBlock.X+2, tl.bottomLeftBlock.Y+y)
	}
	return points
}

type tetrisPlus struct {
	bottomBlock image.Point
}

func makeTetrisPlus(gridHeight int) tetrisShape {
	// offset X by 1 more, we can't be in the bottom left here
	return &tetrisPlus{image.Pt(3, gridHeight+3)}
}

func (tl *tetrisPlus) getHeight() int { return 3 }
func (tl *tetrisPlus) maybeWind(wind int, tg tetrisGrid) {
	for y := 0; y < 3; y++ {
		// X for top and bottom block are equal
		// X for center row is +- 1 from this, in the same direction as the wind
		dest := image.Pt(tl.bottomBlock.X+aElseB(y == 1, wind, 0)+wind, tl.bottomBlock.Y+y)
		if !tg.isOpen(dest) {
			return
		}
	}
	tl.bottomBlock.X += wind
}
func (tl *tetrisPlus) maybeFall(tg tetrisGrid) bool {
	for x := -1; x < 2; x++ {
		// The bottom (center) block needs to check 1 space below it's level
		// Either side needs to check the same level the bottom block is already on
		dest := image.Pt(tl.bottomBlock.X+x, tl.bottomBlock.Y-aElseB(x == 0, 1, 0))
		if !tg.isOpen(dest) {
			return false
		}
	}
	tl.bottomBlock.Y -= 1
	return true
}
func (tl *tetrisPlus) stopFalling() []image.Point {
	points := make([]image.Point, 5)
	for x := -1; x < 2; x++ {
		points[x+1] = image.Pt(tl.bottomBlock.X+x, tl.bottomBlock.Y+1)
	}
	points[3] = tl.bottomBlock
	points[4] = image.Pt(tl.bottomBlock.X, tl.bottomBlock.Y+2)
	return points
}

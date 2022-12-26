package a2022

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"regexp"
	"strings"
)

type wrappingGridPoint struct {
	p                     image.Point
	faceNumber            int
	isOpen                bool
	up, down, left, right *wrappingGridPoint
}

func (p *wrappingGridPoint) setEdge(dir walkDirection, dest *wrappingGridPoint) {
	switch dir {
	case up:
		p.up = dest
	case down:
		p.down = dest
	case left:
		p.left = dest
	case right:
		p.right = dest
	default:
		panic("unhandled setEdge direction")
	}
}

type walkDirection int

func (w walkDirection) String() string {
	switch w {
	case up:
		return "up"
	case down:
		return "down"
	case left:
		return "left"
	case right:
		return "right"
	}
	return "unknown"
}

const (
	right walkDirection = iota
	down
	left
	up
)

func (p wrappingGridPoint) walk(facing walkDirection) image.Point {
	switch facing {
	case right:
		return aElseB(p.right.isOpen, p.right.p, p.p)
	case down:
		return aElseB(p.down.isOpen, p.down.p, p.p)
	case left:
		return aElseB(p.left.isOpen, p.left.p, p.p)
	default:
		return aElseB(p.up.isOpen, p.up.p, p.p)
	}
}

const (
	clockwise        string = "R" // also "down" for voids
	counterClockwise string = "L" // also "up" for voids
)

type edgePair struct{ from, to int }

var foldedEdgeMap = map[edgePair]struct {
	direction string
	magnitude int
}{
	{1, 6}: {clockwise, 2},
	{1, 3}: {counterClockwise, 1},
	{1, 2}: {counterClockwise, 2},
	{1, 4}: {clockwise, 0},
	{2, 1}: {counterClockwise, 2},
	{2, 6}: {counterClockwise, 1},
	{2, 5}: {clockwise, 2},
	{2, 3}: {clockwise, 0},
	{3, 1}: {clockwise, 1},
	{3, 2}: {counterClockwise, 0},
	{3, 5}: {counterClockwise, 1},
	{3, 4}: {clockwise, 0},
	{4, 1}: {clockwise, 0},
	{4, 3}: {counterClockwise, 0},
	{4, 5}: {counterClockwise, 0},
	{4, 6}: {clockwise, 1},
	{5, 4}: {clockwise, 0},
	{5, 3}: {counterClockwise, 1},
	{5, 2}: {counterClockwise, 2},
	{5, 1}: {},
}

func day22(in io.Reader, gridSize int) (int, int) {
	allDirections := regexp.MustCompile(`[LR]`)

	scanner := bufio.NewScanner(in)
	var grid [][]*wrappingGridPoint
	rowCount := 0
	completedFaces := 0
	var faceGrid [][]int
	var rowFaces []int
	var instructions string
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			cells := strings.Split(t, "")
			if rowCount%gridSize == 0 {
				rowFaces = make([]int, len(cells)/gridSize)
				rowFaceCount := 0
				for i := range rowFaces {
					rowFaceCount += aElseB(cells[i*gridSize] == " ", 0, 1)
					rowFaces[i] = aElseB(cells[i*gridSize] == " ", 0, completedFaces+rowFaceCount)
				}
				completedFaces += rowFaceCount
				faceGrid = append(faceGrid, rowFaces)
			}
			row := make([]*wrappingGridPoint, len(cells))
			for x, val := range cells {
				if val == " " {
					row[x] = nil
					continue
				}
				row[x] = &wrappingGridPoint{
					image.Pt(x, rowCount),
					rowFaces[x/gridSize],
					val == ".",
					nil, nil, nil, nil,
				}
			}
			grid = append(grid, row)

			rowCount++
		} else {
			scanner.Scan()
			instructions = scanner.Text()
		}
	}
	faceGrid = func() [][]int {
		maxGridLen := maxInt(mapValue(faceGrid, func(row []int) int { return len(row) })...)
		consistentFaceGrid := make([][]int, len(faceGrid))
		for i := range consistentFaceGrid {
			if len(faceGrid[i]) == maxGridLen {
				consistentFaceGrid[i] = faceGrid[i]
			} else {
				newRow := make([]int, maxGridLen)
				for j := range newRow {
					if j >= len(faceGrid[i]) {
						newRow[j] = 0
					} else {
						newRow[j] = faceGrid[i][j]
					}
				}
				consistentFaceGrid[i] = newRow
			}
		}
		return consistentFaceGrid
	}()

	// link X neighbors
	for _, row := range grid {
		nonEmpty := filter(row, func(p *wrappingGridPoint) bool { return p != nil })
		for i, cell := range nonEmpty {
			if i == 0 {
				cell.left = nonEmpty[len(nonEmpty)-1]
				cell.right = nonEmpty[i+1]
			} else if i == len(nonEmpty)-1 {
				cell.left = nonEmpty[i-1]
				cell.right = nonEmpty[0]
			} else {
				cell.left = nonEmpty[i-1]
				cell.right = nonEmpty[i+1]
			}
		}
	}
	// link Y neighbors
	maxCol := maxInt(mapValue(grid, func(row []*wrappingGridPoint) int { return len(row) })...)
	for x := 0; x < maxCol; x++ {
		nonEmpty := filter(
			mapValue(grid, func(row []*wrappingGridPoint) *wrappingGridPoint {
				if x >= len(row) {
					return nil
				}
				return row[x]
			}),
			func(p *wrappingGridPoint) bool { return p != nil },
		)
		for i, cell := range nonEmpty {
			if i == 0 {
				cell.up = nonEmpty[len(nonEmpty)-1]
				cell.down = nonEmpty[i+1]
			} else if i == len(nonEmpty)-1 {
				cell.up = nonEmpty[i-1]
				cell.down = nonEmpty[0]
			} else {
				cell.up = nonEmpty[i-1]
				cell.down = nonEmpty[i+1]
			}
		}
	}

	var walker image.Point
	for x, cell := range grid[0] {
		if cell == nil {
			continue
		}
		if cell.isOpen {
			walker = image.Pt(x, 0)
			break
		}
	}

	distances := mapValue(allNumbers.FindAllString(instructions, -1), mustInt)
	turns := allDirections.FindAllString(instructions, -1)

	facing := right
	changeFacing := func(dir string) {
		if dir == counterClockwise {
			facing = map[walkDirection]walkDirection{
				right: up,
				down:  right,
				left:  down,
				up:    left,
			}[facing]
		} else {
			facing = map[walkDirection]walkDirection{
				right: down,
				down:  left,
				left:  up,
				up:    right,
			}[facing]
		}
	}
	walkSteps := func(d int) {
		for i := 0; i < d; i++ {
			walker = grid[walker.Y][walker.X].walk(facing)
		}
	}

	for i, steps := range distances {
		walkSteps(steps)
		if i < len(turns) {
			changeFacing(turns[i])
		}
	}
	part1 := (1000 * (walker.Y + 1)) + (4 * (walker.X + 1)) + int(facing)

	type faceEdge struct {
		face int
		edge walkDirection
	}
	//edgeRemapped := map[faceEdge]bool{}
	//for y, row := range faceGrid {
	//	for x, face := range row {
	//		if face == 0 {
	//			continue
	//		}
	//		// init all as false so map is always complete
	//		for _, dir := range []walkDirection{up, down, left, right} {
	//			edgeRemapped[faceEdge{face, dir}] = false
	//		}
	//		if x > 0 {
	//			edgeRemapped[faceEdge{face, left}] = faceGrid[y][x-1] > 0
	//		}
	//		if x < len(faceGrid[y])-1 {
	//			edgeRemapped[faceEdge{face, right}] = faceGrid[y][x+1] > 0
	//		}
	//		if y > 0 {
	//			edgeRemapped[faceEdge{face, up}] = faceGrid[y-1][x] > 0
	//		}
	//		if y < len(faceGrid)-1 {
	//			edgeRemapped[faceEdge{face, down}] = faceGrid[y+1][x-1] > 0
	//		}
	//	}
	//}
	faceTopLeft := map[int]image.Point{}
	for y, row := range faceGrid {
		for x, face := range row {
			if face != 0 {
				faceTopLeft[face] = image.Pt(x*gridSize, y*gridSize)
			}
		}
	}

	// todo: this should be derivable from faceGrid
	magicRemap := map[faceEdge]faceEdge{
		// control input
		//{1, left}:  {3, up},
		//{1, up}:    {2, up},
		//{1, right}: {6, right},
		//{3, down}:  {5, left},
		//{4, right}: {6, up},
		//{2, down}:  {5, down},
		//{2, left}:  {6, down},
		{2, down}:  {3, right},
		{3, left}:  {4, up},
		{5, down}:  {6, right},
		{1, left}:  {4, left},
		{5, right}: {2, right},
		{1, up}:    {6, left},
		{2, up}:    {6, down},
	}
	reverseMagicRemap := map[faceEdge]faceEdge{
		// control input
		//{3, up}:    {1, left},
		//{2, up}:    {1, up},
		//{6, right}: {1, right},
		//{5, left}:  {3, down},
		//{6, up}:    {4, right},
		//{5, down}:  {2, down},
		//{6, down}:  {2, left},
		{3, right}: {2, down},
		{4, up}:    {3, left},
		{6, right}: {5, down},
		{4, left}:  {1, left},
		{2, right}: {5, right},
		{6, left}:  {1, up},
		{6, down}:  {2, up},
	}
	getEdgeFromFace := func(face int, edge walkDirection) []*wrappingGridPoint {
		topLeft := faceTopLeft[face]
		switch edge {
		case up, down:
			return grid[topLeft.Y+aElseB(edge == up, 0, gridSize-1)][topLeft.X : topLeft.X+gridSize]
		default:
			return mapValue(grid[topLeft.Y:topLeft.Y+gridSize], func(row []*wrappingGridPoint) *wrappingGridPoint {
				return row[topLeft.X+aElseB(edge == left, 0, gridSize-1)]
			})
		}
	}
	for from, to := range magicRemap {
		fromPoints := getEdgeFromFace(from.face, from.edge)
		toPoints := getEdgeFromFace(to.face, to.edge)
		if (from.edge == to.edge) ||
			(from.edge == left && to.edge == down) || (from.edge == down && to.edge == left) ||
			(from.edge == right && to.edge == up) || (from.edge == up && to.edge == right) {
			toPoints = reverse(toPoints)
		}
		for i, fromPoint := range fromPoints {
			toPoint := toPoints[i]
			fromPoint.setEdge(from.edge, toPoint)
			toPoint.setEdge(to.edge, fromPoint)
		}
	}

	// reset walker
	facing = right
	for x, cell := range grid[0] {
		if cell == nil {
			continue
		}
		if cell.isOpen {
			walker = image.Pt(x, 0)
			break
		}
	}
	walkStepsCube := func(d int) {
		fromFace := faceEdge{grid[walker.Y][walker.X].faceNumber, facing}
		for i := 0; i < d; i++ {
			walker = grid[walker.Y][walker.X].walk(facing)
			if grid[walker.Y][walker.X].faceNumber != fromFace.face {
				to, hasRemap := magicRemap[fromFace]
				if !hasRemap {
					to, hasRemap = reverseMagicRemap[fromFace]
					if !hasRemap {
						fmt.Printf("walked from face %d to face %d\n", fromFace.face, grid[walker.Y][walker.X].faceNumber)
						fromFace.face = grid[walker.Y][walker.X].faceNumber
						// edges were connected by adjacency
						continue
					}
				}
				facing = (to.edge + 2) % 4
				fmt.Printf("walked from face %d to face %d, now looking %s\n", fromFace.face, to.face, facing)
				fmt.Printf("walker arrived at %v\n", walker.Add(image.Pt(1, 1)))
				fromFace = to
			}
		}
	}
	fmt.Printf("starting from %v facing %s\n", walker, facing)
	// re-walk the grid
	for i, steps := range distances {
		fmt.Printf("walking %d steps\n", steps)
		walkStepsCube(steps)
		fmt.Printf("now at %v\n", walker.Add(image.Pt(1, 1)))
		if i < len(turns) {
			fmt.Printf("turning %s\n", turns[i])
			changeFacing(turns[i])
			fmt.Printf("now facing %s\n", facing)
		}
	}
	//fmt.Printf("ended facing %s (%d) at %v\n", facing, facing, walker)
	part2 := (1000 * (walker.Y + 1)) + (4 * (walker.X + 1)) + int(facing)

	return part1, part2
}

/**
part 2:
if filter(nonEmpty) / 50 > 0
	we have N+1 faces present
assign cubeFace:
	cubeRow = y / 50
	cubeFace = (nonEmpty X / 50 + completed faces)
completedFaces, maxFaces
	rowFaces = (grid[y] / 50) + 1
	if rowFaces != maxFaces
		completedFaces += maxFaces
		maxFaces = rowFaces

folding: yikes. probably some trick here with a mini grid representing the cube faces
	1 opp 5
	3 opp 6
	2 opp 4

if you want <-> into void:
	y%50 -> x offset on destination
	aElseB(newFacing == down, 0, 49) -> y offset on destination

	wrap into same row - e.g 2 walks into the same void above 6 in control
	check above and below - only max one can be occupied, that's where you end up
		iterate until you find a face
			2 walking up: right (or left) twice = 1, going
			2 walking down: right empty, left once = 6, going "up" 6 from the bottom
		rotateDirection := aElseB(
			above,
			aElseB(facing == left, clockwise, counterClockwise),
			aElseB(facing == left, counterClockwise, clockwise),
		)
		i*()
{2,down}:{3,right}
{3,left}:{4,up}
{5,down}:{6,right}
{1,left}:{4,left}
{5,right}:{2,right}
{1,up}:{6,left}
{2,up}:{6,down}
.12
.3.
45.
6..

direct:
1->2
1-v3
3-v5
4->5
4-v6

indirect, 1:
2-v3(r)
4-^3(l)
5-v6(r)

indirect, 2
1->4(l)
2-v5(r)


..1.
234.
..56

direct: 5
close: 3
1->3 (1cc,3c)
3->5 (3cc,5c)
4->6 (4c,6cc)

2->6 (2cc,6c?)

far:4
2->1
2->5
1->6
2->6

*/

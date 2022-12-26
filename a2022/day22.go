package a2022

import (
	"bufio"
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
		// personal input
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
		// personal input
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
						fromFace.face = grid[walker.Y][walker.X].faceNumber
						// edges were connected by adjacency
						continue
					}
				}
				facing = map[walkDirection]walkDirection{
					right: left, left: right,
					up: down, down: up,
				}[to.edge]
				fromFace = to
			}
		}
	}

	// re-walk the grid
	for i, steps := range distances {
		walkStepsCube(steps)
		if i < len(turns) {
			changeFacing(turns[i])
		}
	}

	part2 := (1000 * (walker.Y + 1)) + (4 * (walker.X + 1)) + int(facing)

	return part1, part2
}

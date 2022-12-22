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
	isOpen                bool
	up, down, left, right *wrappingGridPoint
}

func (p wrappingGridPoint) walk(facing int) image.Point {
	switch facing {
	case 0:
		return aElseB(p.right.isOpen, p.right.p, p.p)
	case 1:
		return aElseB(p.down.isOpen, p.down.p, p.p)
	case 2:
		return aElseB(p.left.isOpen, p.left.p, p.p)
	default:
		return aElseB(p.up.isOpen, p.up.p, p.p)
	}
}

var allDirections = regexp.MustCompile(`[LR]`)

func day22(in io.Reader) int {
	scanner := bufio.NewScanner(in)
	var grid [][]*wrappingGridPoint
	rowCount := 0
	var instructions string
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			cells := strings.Split(t, "")
			row := make([]*wrappingGridPoint, len(cells))
			for x, val := range cells {
				if val == " " {
					row[x] = nil
					continue
				}
				row[x] = &wrappingGridPoint{
					image.Pt(x, rowCount),
					val == ".",
					nil, nil, nil, nil,
				}
			}
			grid = append(grid, row)
		} else {
			scanner.Scan()
			instructions = scanner.Text()
		}
		rowCount++
	}
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

	//fmt.Printf("grid lines: %d\ninstructions: %s\n", len(grid), instructions)
	//fmt.Printf("numbers: %v\n", distances)
	//fmt.Printf("turns: %v\n", turns)
	//fmt.Printf("walker starts at: %v (+1 each for final answer)\n", walker)

	facing := 0
	changeFacing := func(dir string) {
		facing += aElseB(dir == "L", -1, 1)
		facing = aElseB(facing%4 >= 0, facing%4, (facing+4)%4)
	}
	walkSteps := func(d int) {
		for i := 0; i < d; i++ {
			walker = grid[walker.Y][walker.X].walk(facing)
		}
	}

	for i, steps := range distances {
		//fmt.Printf("walking %d\n", steps)
		walkSteps(steps)
		if i < len(turns) {
			//fmt.Printf("turning %s\n", turns[i])
			changeFacing(turns[i])
		}
		//fmt.Printf("walkers now facing %d at %v\n", facing, walker)
	}

	//fmt.Printf("walkers ends facing %d at %v\n", facing, walker)

	return (1000 * (walker.Y + 1)) + (4 * (walker.X + 1)) + facing
}

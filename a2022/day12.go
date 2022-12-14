package a2022

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

func letterToInt(letter string) int {
	const alphabetIndex = "_abcdefghijklmnopqrstuvwxyz"
	if letter == "S" {
		return 1
	}
	if letter == "E" {
		return 26
	}
	return strings.Index(alphabetIndex, letter)
}

// used while filling out the node scores, in descending score/search
func getValidVisitors(grid [][]int, p coordinates) []coordinates {
	minNext := grid[p.Y][p.X] - 1

	var neighbors []coordinates
	if p.X > 0 && grid[p.Y][p.X-1] >= minNext {
		neighbors = append(neighbors, coordinates{p.X - 1, p.Y})
	}
	if p.Y > 0 && grid[p.Y-1][p.X] >= minNext {
		neighbors = append(neighbors, coordinates{p.X, p.Y - 1})
	}
	if p.X < len(grid[p.Y])-1 && grid[p.Y][p.X+1] >= minNext {
		neighbors = append(neighbors, coordinates{p.X + 1, p.Y})
	}
	if p.Y < len(grid)-1 && grid[p.Y+1][p.X] >= minNext {
		neighbors = append(neighbors, coordinates{p.X, p.Y + 1})
	}
	return neighbors
}

type pathNode struct {
	pos           coordinates
	height, score int
	neighbors     []*pathNode
}

func (p *pathNode) setScore(score int) {
	p.score = score
}
func (p *pathNode) addNeighbor(n *pathNode) {
	p.neighbors = append(p.neighbors, n)
	sort.SliceStable(p.neighbors, func(i, j int) bool {
		return p.neighbors[i].score > p.neighbors[j].score
	})
}

func day12(in io.Reader) (int, int) {
	s := bufio.NewScanner(in)
	var grid [][]int
	start, end := coordinates{0, 0}, coordinates{0, 0}
	var allStarts []coordinates
	for s.Scan() {
		row := s.Text()
		rowI := len(grid)

		exploded := strings.Split(row, "")
		for i, letter := range exploded {
			switch letter {
			case "S":
				start = coordinates{i, rowI}
				allStarts = append(allStarts, start)
			case "a":
				allStarts = append(allStarts, coordinates{i, rowI})
			case "E":
				end = coordinates{i, rowI}
			}
		}

		grid = append(grid, mapValue(exploded, letterToInt))
	}

	nodeGrid := mapValueWithIndex(grid, func(y int, row []int) []*pathNode {
		return mapValueWithIndex(row, func(x, height int) *pathNode {
			return &pathNode{coordinates{x, y}, height, -1, nil}
		})
	})

	if shouldLog {
		fmt.Printf("start %v, end %v\n", start, end)
	}

	endNode := nodeGrid[end.Y][end.X]
	endNode.setScore(len(grid) * len(grid[0]))
	queue := []*pathNode{endNode}
	var currentQueueNode *pathNode
	for len(queue) > 0 {
		currentQueueNode, queue = queue[0], queue[1:]
		possibleVisitors := getValidVisitors(grid, currentQueueNode.pos)

		if len(possibleVisitors) == 0 {
			continue
		}

		visitorScore := currentQueueNode.score - 1
		for _, v := range possibleVisitors {
			node := nodeGrid[v.Y][v.X]
			node.addNeighbor(currentQueueNode)
			if node.score == -1 {
				node.setScore(visitorScore)
				queue = append(queue, node)
			}
		}
	}

	if shouldLog {
		scoreGrid := mapValue(nodeGrid, func(row []*pathNode) []int { return mapValue(row, func(n *pathNode) int { return n.score }) })
		fmt.Println("score grid")
		for _, row := range scoreGrid {
			fmt.Printf("%v\n", row)
		}

		fmt.Printf("starting coords: %v\n", allStarts)
	}

	max := len(nodeGrid) * len(nodeGrid[0])
	startingSteps, bestSteps := -1, max
	for _, startingCoords := range allStarts {
		steps, currentNode := 0, nodeGrid[startingCoords.Y][startingCoords.X]
		for {
			// If we're stuck, give up - this isn't a winning starting location
			if len(currentNode.neighbors) == 0 {
				steps = max
				break
			}

			steps++
			currentNode = currentNode.neighbors[0]
			if currentNode.pos == end {
				break
			}
		}
		if steps < bestSteps {
			bestSteps = steps
		}
		if startingCoords == start {
			startingSteps = steps
		}
	}
	return startingSteps, bestSteps
}

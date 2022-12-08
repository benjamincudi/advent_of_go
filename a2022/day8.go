package a2022

import (
	"bufio"
	"io"
	"strings"
)

func day8(in io.Reader) (int, int) {
	scanner := bufio.NewScanner(in)

	var forest [][]int
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		forest = append(forest, mapValue(strings.Split(t, ""), func(t string) int {
			return mustInt(t)
		}))
	}
	visible := (len(forest) * 2) + (len(forest[0]) * 2) - 4
	scenicScore := 0
	//fmt.Printf("init visible %d\n", visible)
	for i := 1; i < len(forest)-1; i++ {
		for j := 1; j < len(forest[i])-1; j++ {
			treeHeight := forest[i][j]
			isVisible := false
			// up
			direction := true
			up, down, left, right := 0, 0, 0, 0
			for y := i - 1; y >= 0 && direction; y-- {
				up++
				if treeHeight <= forest[y][j] {
					direction = false
					break
				}
			}
			if direction {
				//fmt.Printf("tree height %d is visible at %d,%d going up\n", treeHeight, i, j)
				isVisible = true
			}

			// down
			direction = true
			for y := i + 1; y < len(forest) && direction; y++ {
				down++
				if treeHeight <= forest[y][j] {
					direction = false
					break
				}
			}
			if direction {
				//fmt.Printf("tree height %d is visible at %d,%d going down\n", treeHeight, i, j)
				isVisible = true
			}

			// left
			direction = true
			for x := j - 1; x >= 0 && direction; x-- {
				left++
				if treeHeight <= forest[i][x] {
					direction = false
					break
				}
			}
			if direction {
				//fmt.Printf("tree height %d is visible at %d,%d going left\n", treeHeight, i, j)
				isVisible = true
			}

			// right
			direction = true
			for x := j + 1; x < len(forest[i]) && direction; x++ {
				right++
				if treeHeight <= forest[i][x] {
					direction = false
					break
				}
			}
			if direction {
				//fmt.Printf("tree height %d is visible at %d,%d going right\n", treeHeight, i, j)
				isVisible = true
			}

			if isVisible {
				visible++
			}
			scenicScore = maxInt(scenicScore, up*left*right*down)
		}
	}
	return visible, scenicScore
}

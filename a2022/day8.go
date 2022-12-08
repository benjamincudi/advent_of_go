package a2022

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type direction struct {
	trees   int
	visible bool
}

func checkDirections(index int, rowOrCol []int) (leftOrUp, rightOrDown direction) {
	treeHeight := rowOrCol[index]
	leftOrUp, rightOrDown = direction{index, true}, direction{(len(rowOrCol) - 1) - index, true}
	for i, height := range rowOrCol {
		if height >= treeHeight {
			if i < index {
				leftOrUp = direction{index - i, false}
			}
			if i > index {
				rightOrDown = direction{i - index, false}
				break
			}
		}
	}
	return leftOrUp, rightOrDown
}

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
	maxDown, maxRight := len(forest)-1, len(forest[0])-1
	if shouldLog {
		fmt.Printf("maxDown %d, maxRight %d\n", maxDown, maxRight)
	}
	visible := (len(forest) * 2) + (len(forest[0]) * 2) - 4
	scenicScore := 0
	if shouldLog {
		fmt.Printf("init visible %d\n", visible)
	}
	for i := 1; i < len(forest)-1; i++ {
		for j := 1; j < len(forest[i])-1; j++ {
			if shouldLog {
				fmt.Printf("h: %d at [%d][%d]\n", forest[i][j], i, j)
			}

			left, right := checkDirections(j, forest[i])
			up, down := checkDirections(i, mapValue(forest, func(row []int) int { return row[j] }))

			anyVisible := up.visible || down.visible || left.visible || right.visible
			score := up.trees * left.trees * right.trees * down.trees
			if shouldLog {
				fmt.Printf("l: %v, r: %v, u: %v, d: %v, anyVisible: %t, score: %d\n", left, right, up, down, anyVisible, score)
			}
			if anyVisible {
				visible++
			}
			scenicScore = maxInt(scenicScore, score)
		}
	}
	return visible, scenicScore
}

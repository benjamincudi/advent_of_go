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
			treeHeight := forest[i][j]
			if shouldLog {
				fmt.Printf("h: %d at [%d][%d]\n", treeHeight, i, j)
			}
			left, right := direction{j, true}, direction{maxRight - j, true}
			for x, height := range forest[i] {
				if height >= treeHeight {
					if x < j {
						left = direction{j - x, false}
					}
					if x > j {
						right = direction{x - j, false}
						break
					}
				}
			}

			up, down := direction{i, true}, direction{maxDown - i, true}
			for y := range forest {
				if forest[y][j] >= treeHeight {
					if y < i {
						up = direction{i - y, false}
					}
					if y > i {
						down = direction{y - i, false}
						break
					}
				}
			}
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

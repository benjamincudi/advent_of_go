package a2024

import (
	"io"
	"strings"
)

func day9(r io.Reader) (int, int) {
	var rawInput []int
	mapFileLines(r, func(text string) {
		rawInput = mapValue(strings.Split(text, ""), mustInt)
	})
	totalLength := 0
	for _, size := range rawInput {
		totalLength += size
	}
	explodedInput := make([]*int, 0, totalLength)
	startIndex := 0
	nilIndexes := make([]int, 0, totalLength)

	type blockInfo struct {
		startIndex, size int
		id               *int
		isFile           bool
	}
	blockInfos := make([]blockInfo, 0, len(rawInput))

	for i, size := range rawInput {
		id := aElseB(i%2 == 0, pointer(i/2), nil)
		blockInfos = append(blockInfos, blockInfo{startIndex, size, id, i%2 == 0})

		blocks := make([]*int, size)
		for j := range blocks {
			if id == nil {
				nilIndexes = append(nilIndexes, startIndex+j)
			}
			blocks[j] = id
		}
		explodedInput = append(explodedInput, blocks...)
		startIndex += size
	}

	copiedResult := make([]*int, len(explodedInput))
	copy(copiedResult, explodedInput)

	var moveTo int
	for i := totalLength - 1; i > nilIndexes[0]; i-- {
		if explodedInput[i] == nil {
			continue
		}
		moveTo, nilIndexes = nilIndexes[0], nilIndexes[1:]
		explodedInput[moveTo] = explodedInput[i]
		explodedInput[i] = nil
	}

	updatedChecksum := 0
	for i, id := range explodedInput {
		if id == nil {
			continue
		}
		updatedChecksum += i * *id
	}

	for i := len(blockInfos) - 1; i > 0; i-- {
		b := blockInfos[i]
		if !b.isFile {
			// only move files
			continue
		}
		for j := 0; j < i; j++ {
			e := blockInfos[j]
			if e.isFile || e.size < b.size {
				// only swap with empty space
				continue
			}
			blockInfos[j].size -= b.size
			blockInfos[j].startIndex += b.size
			for d := 0; d < b.size; d++ {
				copiedResult[e.startIndex+d] = b.id
				copiedResult[b.startIndex+d] = nil
			}
			break
		}
	}

	betterChecksum := 0
	for i, id := range copiedResult {
		if id == nil {
			continue
		}
		betterChecksum += i * *id
	}

	return updatedChecksum, betterChecksum
}

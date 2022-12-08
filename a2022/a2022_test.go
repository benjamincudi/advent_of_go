package a2022

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"
)

type failerHelper interface {
	Fatalf(format string, args ...any)
	Helper()
}

func mustOpen(t failerHelper, name string) fs.File {
	t.Helper()
	inputReader, err := inputs.Open(fmt.Sprintf("inputs-2022/%s", name))
	if err != nil {
		t.Fatalf("%v", err)
	}
	return inputReader
}

func Test_day1(t *testing.T) {
	testCases := []struct {
		name       string
		in         io.Reader
		top1, top3 int
	}{
		{"control", mustOpen(t, "control1.txt"), 24000, 45000},
		{"personal", mustOpen(t, "day1.txt"), 71924, 210406},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			top3 := day1(tc.in)
			sum := 0
			for _, cal := range top3 {
				sum += cal
			}
			if top3[0] != tc.top1 {
				t.Errorf("expected max %d, got %d", tc.top1, top3[0])
			}
			if sum != tc.top3 {
				t.Errorf("expected top 3 sum of %d, got %d", tc.top3, sum)
			}
		})
	}
}

func Test_day2(t *testing.T) {
	testCases := []struct {
		name                string
		in                  io.Reader
		secondColumnAsPlays bool
		outcome             int
	}{
		{"control case - as plays", mustOpen(t, "control2.txt"), true, 15},
		{"control case - as outcomes", mustOpen(t, "control2.txt"), false, 12},
		{"personal input - as plays", mustOpen(t, "day2.txt"), true, 13565},
		{"personal input - as outcomes", mustOpen(t, "day2.txt"), false, 12424},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if v := day2(tc.in, tc.secondColumnAsPlays); v != tc.outcome {
				t.Errorf("expected %d, got %d", tc.outcome, v)
			}
		})
	}
}

func Test_day3(t *testing.T) {
	testCases := []struct {
		name         string
		in           io.Reader
		sacks, group int
	}{
		{"control case", mustOpen(t, "control3.txt"), 157, 70},
		{"personal input", mustOpen(t, "day3.txt"), 7967, 2716},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sacks, group, err := day3(tc.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if sacks != tc.sacks {
				t.Errorf("expected %d, got %d", tc.sacks, sacks)
			}
			if group != tc.group {
				t.Errorf("expected group score %d, got %d", tc.group, group)
			}
		})
	}
}

func Test_day4(t *testing.T) {
	testCases := []struct {
		name                              string
		in                                io.Reader
		fullOverlapCount, anyOverlapCount int
	}{
		{"control case", mustOpen(t, "control4.txt"), 2, 4},
		{"personal input", mustOpen(t, "day4.txt"), 441, 861},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if fullOverlapCount, anyOverlapCount, err := day4(tc.in); err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if fullOverlapCount != tc.fullOverlapCount {
				t.Errorf("expected %d overlapping fullOverlap, got %d", tc.fullOverlapCount, fullOverlapCount)
			} else if anyOverlapCount != tc.anyOverlapCount {
				t.Errorf("expected %d with any overlap, got %d", tc.anyOverlapCount, anyOverlapCount)
			}
		})
	}
}

func Test_day5(t *testing.T) {
	testCases := []struct {
		name         string
		in           io.Reader
		part1, part2 string
	}{
		{"control case", mustOpen(t, "control5.txt"), "CMZ", "MCD"},
		{"personal input", mustOpen(t, "day5.txt"), "FZCMJCRHZ", "JSDHQMZGF"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2 := day5(tc.in)
			if part1 != tc.part1 {
				t.Errorf("expected part1 to be %s, got %s", tc.part1, part1)
			}
			if part2 != tc.part2 {
				t.Errorf("expected part2 to be %s, got %s", tc.part2, part2)
			}
		})
	}
}

func Test_day6(t *testing.T) {
	reader := bufio.NewScanner(mustOpen(t, "control6.txt"))
	var control []string
	for reader.Scan() {
		control = append(control, reader.Text())
	}
	testCases := []struct {
		name                          string
		in                            io.Reader
		startOfPacket, startOfMessage int
	}{
		{"control case 1", strings.NewReader(control[0]), 7, 19},
		{"control case 2", strings.NewReader(control[1]), 5, 23},
		{"control case 3", strings.NewReader(control[2]), 6, 23},
		{"control case 4", strings.NewReader(control[3]), 10, 29},
		{"control case 5", strings.NewReader(control[4]), 11, 26},
		{"personal input", mustOpen(t, "day6.txt"), 1598, 2414},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := day6(tc.in); r.startOfPacket != tc.startOfPacket {
				t.Errorf("expected start-of-packet %d, got %d", tc.startOfPacket, r.startOfPacket)
			} else if r.startOfMessage != tc.startOfMessage {
				t.Errorf("expected start-of-message %d, got %d", tc.startOfMessage, r.startOfMessage)
			}
		})
	}
	testCases = []struct {
		name                          string
		in                            io.Reader
		startOfPacket, startOfMessage int
	}{
		{"control case 1", strings.NewReader(control[0]), 7, 19},
		{"control case 2", strings.NewReader(control[1]), 5, 23},
		{"control case 3", strings.NewReader(control[2]), 6, 23},
		{"control case 4", strings.NewReader(control[3]), 10, 29},
		{"control case 5", strings.NewReader(control[4]), 11, 26},
		{"personal input", mustOpen(t, "day6.txt"), 1598, 2414},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("fast - %s", tc.name), func(t *testing.T) {
			if r := day6fast(tc.in); r.startOfPacket != tc.startOfPacket {
				t.Errorf("expected start-of-packet %d, got %d", tc.startOfPacket, r.startOfPacket)
			} else if r.startOfMessage != tc.startOfMessage {
				t.Errorf("expected start-of-message %d, got %d", tc.startOfMessage, r.startOfMessage)
			}
		})
	}
}

func Benchmark_day6(b *testing.B) {
	reader := bufio.NewScanner(mustOpen(b, "control6.txt"))
	var testCases []string
	for reader.Scan() {
		testCases = append(testCases, reader.Text())
	}
	reader = bufio.NewScanner(mustOpen(b, "day6.txt"))
	for reader.Scan() {
		testCases = append(testCases, reader.Text())
	}
	modulo := len(testCases)
	b.Run("day6 normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			day6(strings.NewReader(testCases[i%modulo]))
		}
	})
	b.Run("day6 fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			day6fast(strings.NewReader(testCases[i%modulo]))
		}
	})
}

func Test_day7(t *testing.T) {
	testCases := []struct {
		name                         string
		in                           io.Reader
		sumOfSmall, smallestToDelete int
	}{
		{"control case", mustOpen(t, "control7.txt"), 95437, 24933642},
		{"personal input", mustOpen(t, "day7.txt"), 1453349, 2948823},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if sum, small := day7(tc.in); sum != tc.sumOfSmall {
				t.Errorf("expected sumOfSmall %d, got %d", tc.sumOfSmall, sum)
			} else if small != tc.smallestToDelete {
				t.Errorf("expected smallestToDelete %d, got %d", tc.smallestToDelete, small)
			}
		})
	}
}

func Test_day8(t *testing.T) {
	testCases := []struct {
		name           string
		in             io.Reader
		visible, score int
	}{
		{"control case", mustOpen(t, "control8.txt"), 21, 8},
		{"personal input", mustOpen(t, "day8.txt"), 1679, 536625},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visible, score := day8(tc.in)
			if visible != tc.visible {
				t.Errorf("expected visible %v, got %v", tc.visible, visible)
			}
			if score != tc.score {
				t.Errorf("expected score %d, got %d", tc.score, score)
			}
		})
	}
}

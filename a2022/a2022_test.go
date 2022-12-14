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
		{"control", mustOpen(t, "control01.txt"), 24000, 45000},
		{"personal", mustOpen(t, "day01.txt"), 71924, 210406},
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
		{"control case - as plays", mustOpen(t, "control02.txt"), true, 15},
		{"control case - as outcomes", mustOpen(t, "control02.txt"), false, 12},
		{"personal input - as plays", mustOpen(t, "day02.txt"), true, 13565},
		{"personal input - as outcomes", mustOpen(t, "day02.txt"), false, 12424},
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
		{"control case", mustOpen(t, "control03.txt"), 157, 70},
		{"personal input", mustOpen(t, "day03.txt"), 7967, 2716},
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
		{"control case", mustOpen(t, "control04.txt"), 2, 4},
		{"personal input", mustOpen(t, "day04.txt"), 441, 861},
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
		{"control case", mustOpen(t, "control05.txt"), "CMZ", "MCD"},
		{"personal input", mustOpen(t, "day05.txt"), "FZCMJCRHZ", "JSDHQMZGF"},
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
	reader := bufio.NewScanner(mustOpen(t, "control06.txt"))
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
		{"personal input", mustOpen(t, "day06.txt"), 1598, 2414},
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
		{"personal input", mustOpen(t, "day06.txt"), 1598, 2414},
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
	reader := bufio.NewScanner(mustOpen(b, "control06.txt"))
	var testCases []string
	for reader.Scan() {
		testCases = append(testCases, reader.Text())
	}
	reader = bufio.NewScanner(mustOpen(b, "day06.txt"))
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
		{"control case", mustOpen(t, "control07.txt"), 95437, 24933642},
		{"personal input", mustOpen(t, "day07.txt"), 1453349, 2948823},
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
		{"control case", mustOpen(t, "control08.txt"), 21, 8},
		{"personal input", mustOpen(t, "day08.txt"), 1679, 536625},
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

func Test_day9(t *testing.T) {
	testCases := []struct {
		name                          string
		in                            io.Reader
		tailVisited, chainTailVisited int
	}{
		{"control case", mustOpen(t, "control09.txt"), 13, 1},
		{"control case 2", mustOpen(t, "control09-2.txt"), 88, 36},
		{"personal input", mustOpen(t, "day09.txt"), 6494, 2691},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tailVisited, chainTailVisited := day9(tc.in)
			if tailVisited != tc.tailVisited {
				t.Errorf("expected tailVisited %v, got %v", tc.tailVisited, tailVisited)
			}
			if chainTailVisited != tc.chainTailVisited {
				t.Errorf("expected chainTailVisited %d, got %d", tc.chainTailVisited, chainTailVisited)
			}
		})
	}
}

func Test_day10(t *testing.T) {
	// cycle off-by-one is probably the crux of everyone's bugs
	// if someone is struggling with this in part 2, the screen looks like this:
	// ##..##..#..##...##.##..##..##..##..##..#
	// ###..####..###...###...####..####..###..
	// ##.....#####...###.....####.....##.....#
	// ####.....##.#......#####.....####.......
	// ###.##.....##.###......######......##.#.
	// ###.##.......####.#........########.....
	controlOutputScreen := `
##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....
`
	personalOutputScreen := `
###..#..#..##..#..#.#..#.###..####.#..#.
#..#.#..#.#..#.#.#..#..#.#..#.#....#.#..
#..#.#..#.#..#.##...####.###..###..##...
###..#..#.####.#.#..#..#.#..#.#....#.#..
#.#..#..#.#..#.#.#..#..#.#..#.#....#.#..
#..#..##..#..#.#..#.#..#.###..####.#..#.
`
	testCases := []struct {
		name   string
		in     io.Reader
		sum    int
		screen string
	}{
		{"control case", mustOpen(t, "control10.txt"), 13140, controlOutputScreen},
		{"personal input", mustOpen(t, "day10.txt"), 13220, personalOutputScreen},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sum, screen, err := day10(tc.in)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if sum != tc.sum {
				t.Errorf("expected %v, got %v", tc.sum, sum)
			}
			if screen != tc.screen {
				t.Errorf("expected screen %s\ngot screen %s", tc.screen, screen)
			}
		})
	}
}

func Test_day11(t *testing.T) {
	testCases := []struct {
		name                              string
		in                                io.Reader
		monkeyBusiness, bigMonkeyBusiness int
	}{
		{"control case", mustOpen(t, "control11.txt"), 10605, 2713310158},
		{"personal input", mustOpen(t, "day11.txt"), 58322, 13937702909},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			monkeyBusiness, bigMonkeyBusiness := day11(tc.in)
			if monkeyBusiness != tc.monkeyBusiness {
				t.Errorf("expected monkeyBusiness %v, got %v", tc.monkeyBusiness, monkeyBusiness)
			}
			if bigMonkeyBusiness != tc.bigMonkeyBusiness {
				t.Errorf("expected bigMonkeyBusiness %v, got %v", tc.bigMonkeyBusiness, bigMonkeyBusiness)
			}
		})
	}
}

func Test_day12(t *testing.T) {
	testCases := []struct {
		name               string
		in                 io.Reader
		steps, scenicSteps int
	}{
		{"control case", mustOpen(t, "control12.txt"), 31, 29},
		{"personal input", mustOpen(t, "day12.txt"), 412, 402},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			steps, scenicSteps := day12(tc.in)
			if steps != tc.steps {
				t.Errorf("expected steps %v, got %v", tc.steps, steps)
			}
			if scenicSteps != tc.scenicSteps {
				t.Errorf("expected scenicSteps %v, got %v", tc.scenicSteps, scenicSteps)
			}
		})
	}
}

func Test_day13(t *testing.T) {
	testCases := []struct {
		name         string
		in           io.Reader
		sum, decoder int
	}{
		{"control case", mustOpen(t, "control13.txt"), 13, 140},
		{"personal input", mustOpen(t, "day13.txt"), 5843, 26289},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sum, decoder := day13(tc.in)
			if sum != tc.sum {
				t.Errorf("expected sum %v, got %v", tc.sum, sum)
			}
			if decoder != tc.decoder {
				t.Errorf("expected decoder %v, got %v", tc.decoder, decoder)
			}
		})
	}
}

func Test_day14(t *testing.T) {
	testCases := []struct {
		name                 string
		in                   io.Reader
		sandUnits, withFloor int
	}{
		{"control case", mustOpen(t, "control14.txt"), 24, 93},
		{"personal input", mustOpen(t, "day14.txt"), 825, 26729},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sandUnits, withFloor := day14(tc.in)
			if sandUnits != tc.sandUnits {
				t.Errorf("expected sandUnits %v, got %v", tc.sandUnits, sandUnits)
			}
			if withFloor != tc.withFloor {
				t.Errorf("expected withFloor %v, got %v", tc.withFloor, withFloor)
			}
		})
	}
}

func Test_day15(t *testing.T) {
	testCases := []struct {
		name                               string
		in                                 io.Reader
		targetRow, upperBounds             int
		occupiedPositions, tuningFrequency int
	}{
		{"control case", mustOpen(t, "control15.txt"), 10, 20, 26, 56000011},
		{"personal input", mustOpen(t, "day15.txt"), 2000000, 4000000, 5112034, 13172087230812},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			occupiedPositions, tuningFrequency := day15(tc.in, tc.targetRow, tc.upperBounds)
			if occupiedPositions != tc.occupiedPositions {
				t.Errorf("expected occupiedPositions %v, got %v", tc.occupiedPositions, occupiedPositions)
			}
			if tuningFrequency != tc.tuningFrequency {
				t.Errorf("expected tuningFrequency %v, got %v", tc.tuningFrequency, tuningFrequency)
			}
		})
	}
}

func Test_day16(t *testing.T) {
	testCases := []struct {
		name                           string
		in                             io.Reader
		pressureReleased, withElephant int
	}{
		{"control case", mustOpen(t, "control16.txt"), 1651, 1707},
		{"personal input", mustOpen(t, "day16.txt"), 1857, 2536},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pressureReleased, withElephant := day16(tc.in)
			if pressureReleased != tc.pressureReleased {
				t.Errorf("expected pressureReleased %v, got %v", tc.pressureReleased, pressureReleased)
			}
			if withElephant != tc.withElephant {
				t.Errorf("expected withElephant %v, got %v", tc.withElephant, withElephant)
			}
		})
	}
}

func Test_day17(t *testing.T) {
	testCases := []struct {
		name                   string
		in                     io.Reader
		height, trillionHeight int
	}{
		{"control case", mustOpen(t, "control17.txt"), 3068, 1514285714288},
		{"personal input", mustOpen(t, "day17.txt"), 3106, 1537175792495},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			height, trillionHeight := day17(tc.in)
			if height != tc.height {
				t.Errorf("expected height %v, got %v", tc.height, height)
			}
			if trillionHeight != tc.trillionHeight {
				t.Errorf("expected trillionHeight %v, got %v", tc.trillionHeight, trillionHeight)
			}
		})
	}
}

func Test_day18(t *testing.T) {
	testCases := []struct {
		name                         string
		in                           io.Reader
		surfaceArea, exteriorSurface int
	}{
		{"control case", mustOpen(t, "control18.txt"), 64, 58},
		{"personal input", mustOpen(t, "day18.txt"), 3432, 2042},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			surfaceArea, exteriorSurface := day18(tc.in)
			if surfaceArea != tc.surfaceArea {
				t.Errorf("expected surfaceArea %v, got %v", tc.surfaceArea, surfaceArea)
			}
			if exteriorSurface != tc.exteriorSurface {
				t.Errorf("expected exteriorSurface %v, got %v", tc.exteriorSurface, exteriorSurface)
			}
		})
	}
}

func Test_day19(t *testing.T) {
	testCases := []struct {
		name                          string
		in                            io.Reader
		totalQuality, productQuantity int
	}{
		{"control case", mustOpen(t, "control19.txt"), 33, 3472},
		{"personal input", mustOpen(t, "day19.txt"), 1147, 3080},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totalQuality, productQuantity := day19(tc.in)
			if totalQuality != tc.totalQuality {
				t.Errorf("expected totalQuality %v, got %v", tc.totalQuality, totalQuality)
			}
			if productQuantity != tc.productQuantity {
				t.Errorf("expected productQuantity %v, got %v", tc.productQuantity, productQuantity)
			}
		})
	}
}

func Test_day20(t *testing.T) {
	testCases := []struct {
		name                  string
		in                    io.Reader
		coordSum, withDecrypt int
	}{
		{"control case", mustOpen(t, "control20.txt"), 3, 1623178306},
		{"personal input", mustOpen(t, "day20.txt"), 7225, 548634267428},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coordSum, withDecrypt := day20_ll(tc.in)
			if coordSum != tc.coordSum {
				t.Errorf("expected coordSum %v, got %v", tc.coordSum, coordSum)
			}
			if withDecrypt != tc.withDecrypt {
				t.Errorf("expected withDecrypt %v, got %v", tc.withDecrypt, withDecrypt)
			}
		})
	}
}

func Test_day21(t *testing.T) {
	testCases := []struct {
		name               string
		in                 io.Reader
		rootYell, humnYell int
	}{
		{"control case", mustOpen(t, "control21.txt"), 152, 301},
		{"personal input", mustOpen(t, "day21.txt"), 142707821472432, 3587647562851},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootYell, humnYell := day21(tc.in)
			if rootYell != tc.rootYell {
				t.Errorf("expected rootYell %v, got %v", tc.rootYell, rootYell)
			}
			if humnYell != tc.humnYell {
				t.Errorf("expected humnYell %v, got %v", tc.humnYell, humnYell)
			}
		})
	}
}

func Test_day22(t *testing.T) {
	// control input
	controlEdgeRemap := map[faceEdge]faceEdge{
		{1, left}:  {3, up},
		{1, up}:    {2, up},
		{1, right}: {6, right},
		{3, down}:  {5, left},
		{4, right}: {6, up},
		{2, down}:  {5, down},
		{2, left}:  {6, down},
	}
	// personal input
	personalEdgeRemap := map[faceEdge]faceEdge{
		{2, down}:  {3, right},
		{3, left}:  {4, up},
		{5, down}:  {6, right},
		{1, left}:  {4, left},
		{5, right}: {2, right},
		{1, up}:    {6, left},
		{2, up}:    {6, down},
	}

	testCases := []struct {
		name               string
		in                 io.Reader
		cubeEdgeRemapper   map[faceEdge]faceEdge
		gridSize, password int
		cubePassword       int
	}{
		{"control case", mustOpen(t, "control22.txt"), controlEdgeRemap, 4, 6032, 5031},
		{"personal input", mustOpen(t, "day22.txt"), personalEdgeRemap, 50, 93226, 37415},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			password, cubePassword := day22(tc.in, tc.gridSize, tc.cubeEdgeRemapper)
			if password != tc.password {
				t.Errorf("expected password %v, got %v", tc.password, password)
			}
			if cubePassword != tc.cubePassword {
				t.Errorf("expected cubePassword %v, got %v", tc.cubePassword, cubePassword)
			}
		})
	}
}

func Test_day23(t *testing.T) {
	testCases := []struct {
		name                    string
		in                      io.Reader
		emptySpaces, finalRound int
	}{
		{"control case", mustOpen(t, "control23.txt"), 110, 20},
		{"personal input", mustOpen(t, "day23.txt"), 4034, 960},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emptySpaces, finalRound := day23(tc.in)
			if emptySpaces != tc.emptySpaces {
				t.Errorf("expected emptySpaces %v, got %v", tc.emptySpaces, emptySpaces)
			}
			if finalRound != tc.finalRound {
				t.Errorf("expected finalRound %v, got %v", tc.finalRound, finalRound)
			}
		})
	}
}

func Test_day24(t *testing.T) {
	testCases := []struct {
		name                 string
		in                   io.Reader
		steps, extendedSteps int
	}{
		{"control case", mustOpen(t, "control24.txt"), 18, 54},
		{"personal input", mustOpen(t, "day24.txt"), 269, 825},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			steps, extendedSteps := day24(tc.in)
			if steps != tc.steps {
				t.Errorf("expected steps %v, got %v", tc.steps, steps)
			}
			if extendedSteps != tc.extendedSteps {
				t.Errorf("expected extendedSteps %v, got %v", tc.extendedSteps, extendedSteps)
			}
		})
	}
}

func Test_day25(t *testing.T) {
	testCases := []struct {
		name     string
		in       io.Reader
		snafuSum string
	}{
		{"control case", mustOpen(t, "control25.txt"), "2=-1=0"},
		{"personal input", mustOpen(t, "day25.txt"), "2-0==21--=0==2201==2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			snafuSum := day25(tc.in)
			if snafuSum != tc.snafuSum {
				t.Errorf("expected snafuSum %v, got %v", tc.snafuSum, snafuSum)
			}
		})
	}
}

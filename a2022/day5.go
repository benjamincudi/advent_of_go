package a2022

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"

	"github.com/gocarina/gocsv"
)

var (
	instructionNumbers = regexp.MustCompile("\\d+")                // rip numbers out of "move X from Y to Z"
	crates             = regexp.MustCompile("(\\[\\w]|\\s{3})\\s") // split rows of text into chunks of crates
	emptyCrate         = regexp.MustCompile("\\s+")                // check for all whitespace, which in this usage means "no crate"
	innerCrate         = regexp.MustCompile("\\[(?P<value>\\w)]")  // pull the crate letter out of a match
	columnLabelRegex   = regexp.MustCompile("((\\s\\d\\s)\\s?)+")  // match the row that labels each crate stack with a number
)

type crateStack struct {
	crates []string
}

func (c *crateStack) Add(crates ...string) {
	c.crates = append(
		// copy the crates to avoid slice reference shenanigans
		append(make([]string, 0, len(crates)+len(c.crates)), crates...),
		c.crates...)
}
func (c *crateStack) Take(num int) []string {
	crate := c.crates[0:num]
	c.crates = append([]string{}, c.crates[num:]...)
	return crate
}
func (c crateStack) Top() string {
	return c.crates[0]
}

type instruction struct {
	count, from, to int
}

func (i *instruction) UnmarshalCSV(s string) error {
	matches := instructionNumbers.FindAllString(s, -1)
	if len(matches) != 3 {
		return errors.New(fmt.Sprintf("did not match 3 numbers in %s", s))
	}
	i.count, i.from, i.to = mustInt(matches[0]), mustInt(matches[1])-1, mustInt(matches[2])-1
	return nil
}

func initStacks(crateRows [][]string) []crateStack {
	crateStacks := make([]crateStack, len(crateRows[0]))
	for _, row := range crateRows {
		for i, crate := range row {
			if crate != "-" {
				crateStacks[i].Add(crate)
			}
		}
	}
	return crateStacks
}

func day5(in io.Reader) (string, string) {
	r := bufio.NewReader(in)
	var crateRows [][]string
	for s, err := r.ReadString('\n'); !columnLabelRegex.MatchString(s) && err != io.EOF; s, err = r.ReadString('\n') {
		if err != nil {
			panic(err)
		}
		var row []string
		for _, matchGroup := range crates.FindAllStringSubmatch(s, -1) {
			if emptyCrate.MatchString(matchGroup[1]) {
				row = append(row, "-")
				continue
			}
			row = append(row, innerCrate.FindStringSubmatch(matchGroup[1])[1])
		}
		crateRows = append(crateRows, row)
	}

	// Reverse the rows so we're building from bottom-up
	sort.SliceStable(crateRows, func(i, j int) bool {
		return i > j
	})
	// generally pointless struct so that we can unmarshal instructions
	type insRow struct {
		I instruction
	}
	var ins []insRow
	if err := gocsv.UnmarshalWithoutHeaders(r, &ins); err != nil {
		panic(err)
	}
	crateStacks := initStacks(crateRows)
	for _, row := range ins {
		for i := 0; i < row.I.count; i++ {
			crateStacks[row.I.to].Add(crateStacks[row.I.from].Take(1)...)
		}
	}
	part1 := ""
	for _, c := range crateStacks {
		part1 += c.Top()
	}

	crateStacks = initStacks(crateRows)
	for _, row := range ins {
		crateStacks[row.I.to].Add(crateStacks[row.I.from].Take(row.I.count)...)
	}
	part2 := ""
	for _, c := range crateStacks {
		part2 += c.Top()
	}
	return part1, part2
}

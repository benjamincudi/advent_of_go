package a2022

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/gocarina/gocsv"
)

var (
	instructionNumbers = regexp.MustCompile(`\d+`)              // rip numbers out of "move X from Y to Z"
	crates             = regexp.MustCompile(`(\[\w]|\s{3})\s`)  // split rows of text into chunks of crates
	emptyCrate         = regexp.MustCompile(`\s+`)              // check for all whitespace, which in this usage means "no crate"
	innerCrate         = regexp.MustCompile(`\[(?P<value>\w)]`) // pull the crate letter out of a match
	columnLabelRegex   = regexp.MustCompile(`((\s\d\s)\s?)+`)   // match the row that labels each crate stack with a number
)

type crateStack struct {
	crates []string
}

func (c *crateStack) Add(crates ...string) {
	c.crates = append(crates, c.crates...)
}
func (c *crateStack) Take(num int) []string {
	crate := c.crates[0:num]
	c.crates = append(make([]string, 0, len(c.crates)-num), c.crates[num:]...)
	return crate
}
func (c *crateStack) Top() string {
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

func day5(in io.Reader) (part1, part2 string) {
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
	reverse(crateRows)
	// generally pointless struct so that we can unmarshal instructions
	type insRow struct{ I instruction }
	var ins []insRow
	if err := gocsv.UnmarshalWithoutHeaders(r, &ins); err != nil {
		panic(err)
	}
	crateStacks1, crateStacks2 := initStacks(crateRows), initStacks(crateRows)
	for _, row := range ins {
		// "one at a time" ordering is just the reverse of all-at-once
		crateStacks1[row.I.to].Add(reverse(crateStacks1[row.I.from].Take(row.I.count))...)
		crateStacks2[row.I.to].Add(crateStacks2[row.I.from].Take(row.I.count)...)
	}
	for i := 0; i < len(crateStacks1); i++ {
		part1 += crateStacks1[i].Top()
		part2 += crateStacks2[i].Top()
	}

	return part1, part2
}

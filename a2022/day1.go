package a2022

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gocarina/gocsv"
)

type preserveEmptyLineReader struct {
	r *bufio.Reader
}

func (p preserveEmptyLineReader) Read() ([]string, error) {
	var snacks []string
	var s string
	var err error
	for s, err = p.r.ReadString('\n'); err != io.EOF; s, err = p.r.ReadString('\n') {
		if err != nil {
			return nil, err
		}
		// Empty line means we should stop reading, we're done with current elf
		if s == "\n" {
			break
		}
		snacks = append(snacks, strings.TrimSuffix(s, "\n"))
	}
	return []string{fmt.Sprintf("[%s]", strings.Join(snacks, ","))}, err
}
func (p preserveEmptyLineReader) ReadAll() ([][]string, error) {
	var all [][]string
	for s, err := p.Read(); err != io.EOF; {
		if err != nil {
			return nil, err
		}
		all = append(all, s)
		s, err = p.Read()
	}
	return all, nil
}

func day1(r io.Reader) []int {
	type Calories struct {
		Cal []int
	}

	var elves []Calories

	if err := gocsv.UnmarshalCSVWithoutHeaders(preserveEmptyLineReader{bufio.NewReader(r)}, &elves); err != nil {
		panic(err)
	}

	var elfTotals []int
	for _, elf := range elves {
		currentTotal := 0
		for _, snack := range elf.Cal {
			currentTotal += snack
		}
		elfTotals = append(elfTotals, currentTotal)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(elfTotals)))
	// max: 71924, totalTop3: 210406
	return elfTotals[0:3]
}

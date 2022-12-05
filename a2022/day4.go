package a2022

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
)

type intRange struct {
	min, max int
}

func (i *intRange) contains(j intRange) bool {
	return i.min <= j.min && i.max >= j.max
}

func (i *intRange) overlaps(j intRange) bool {
	return i.min <= j.min && j.min <= i.max
}

func (i *intRange) UnmarshalCSV(s string) error {
	bounds := strings.Split(s, "-")
	if len(bounds) != 2 {
		return errors.New(fmt.Sprintf("unexpected range input: %s", s))
	}
	i.min, i.max = mustInt(bounds[0]), mustInt(bounds[1])
	return nil
}

type elfPairs struct {
	A, B intRange
}

var boolToInt = map[bool]int{true: 1}

func (p elfPairs) fullyOverlap() int {
	return boolToInt[p.A.contains(p.B) || p.B.contains(p.A)]
}

func (p elfPairs) anyOverlap() int {
	return boolToInt[p.A.overlaps(p.B) || p.B.overlaps(p.A)]
}

func day4(in io.Reader) (int, int, error) {
	var pairs []elfPairs
	if err := gocsv.UnmarshalWithoutHeaders(in, &pairs); err != nil {
		return 0, 0, err
	}
	fullyOverlap, anyOverlap := 0, 0
	for _, pair := range pairs {
		fullyOverlap += pair.fullyOverlap()
		anyOverlap += pair.anyOverlap()
	}
	return fullyOverlap, anyOverlap, nil
}

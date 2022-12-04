package a2022

import (
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"strconv"
	"strings"
)

type IntRange struct {
	min, max int
}

func (i *IntRange) contains(j IntRange) bool {
	return i.min <= j.min && i.max >= j.max
}

func (i *IntRange) overlaps(j IntRange) bool {
	return i.min <= j.min && j.min <= i.max
}

func (i *IntRange) UnmarshalCSV(s string) error {
	bounds := strings.Split(s, "-")
	if len(bounds) != 2 {
		return errors.New(fmt.Sprintf("unexpected range input: %s", s))
	}
	min, err := strconv.Atoi(bounds[0])
	if err != nil {
		return err
	}
	max, err := strconv.Atoi(bounds[1])
	if err != nil {
		return err
	}
	i.min, i.max = min, max
	return nil
}

type elfPairs struct {
	A, B IntRange
}

func (p elfPairs) fullyOverlap() bool {
	return p.A.contains(p.B) || p.B.contains(p.A)
}

func (p elfPairs) anyOverlap() bool {
	return p.A.overlaps(p.B) || p.B.overlaps(p.A)
}

func day4(in io.Reader) (int, int, error) {
	var pairs []elfPairs
	if err := gocsv.UnmarshalWithoutHeaders(in, &pairs); err != nil {
		return 0, 0, err
	}
	fullyOverlap := 0
	anyOverlap := 0
	for _, pair := range pairs {
		if pair.fullyOverlap() {
			fullyOverlap += 1
		}
		if pair.anyOverlap() {
			anyOverlap += 1
		}
	}
	return fullyOverlap, anyOverlap, nil
}

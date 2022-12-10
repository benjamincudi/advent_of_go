package a2022

import (
	"errors"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
)

type rucksack struct {
	c1, c2 []string
}

func (r *rucksack) UnmarshalCSV(s string) error {
	contents := strings.Split(s, "")
	if len(contents)%2 != 0 {
		return errors.New("invalid rucksack has odd-count of items")
	}
	half := len(contents) / 2
	r.c1, r.c2 = contents[0:half], contents[half:]
	return nil
}

func (r *rucksack) allItems() []string {
	return append(r.c1, r.c2...)
}

func (r *rucksack) getOverlappingItems() map[string]int {
	contentsMap := func(c []string) map[string]int {
		m := map[string]int{}
		for _, letter := range c {
			m[letter] = m[letter] + 1
		}
		return m
	}

	overlap := map[string]int{}
	c1 := contentsMap(r.c1)
	c2 := contentsMap(r.c2)
	for match, count1 := range c1 {
		if count2, ok := c2[match]; ok {
			overlap[match] = count1 + count2
		}
	}
	return overlap
}

func day3(in io.Reader) (int, int, error) {
	type sack struct {
		R rucksack
	}
	var sacks []sack
	if err := gocsv.UnmarshalWithoutHeaders(in, &sacks); err != nil {
		return 0, 0, err
	}

	intersectOverlaps := func(o1, o2 map[string]int) map[string]int {
		overlap := map[string]int{}
		for match, count1 := range o1 {
			if count2, ok := o2[match]; ok {
				overlap[match] = count1 + count2
			}
		}
		return overlap
	}

	scoreOverlap := func(o map[string]int) int {
		// _ takes the 0 index so that index is directly score
		values := "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		score := 0
		for letter := range o {
			score += strings.Index(values, letter)
		}
		return score
	}

	sackScore := 0
	for _, s := range sacks {
		sackScore += scoreOverlap(s.R.getOverlappingItems())
	}

	groupScore := 0
	for i := 0; i < len(sacks); i += 3 {
		fauxSack1 := rucksack{sacks[i].R.allItems(), sacks[i+1].R.allItems()}
		fauxSack2 := rucksack{sacks[i+1].R.allItems(), sacks[i+2].R.allItems()}
		groupScore += scoreOverlap(intersectOverlaps(fauxSack1.getOverlappingItems(), fauxSack2.getOverlappingItems()))
	}
	return sackScore, groupScore, nil
}

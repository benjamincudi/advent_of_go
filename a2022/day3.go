package a2022

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"strings"
)

type rucksack struct {
	c1, c2 []string
}

func (r rucksack) allItems() []string {
	all := append(r.c1, r.c2...)
	sort.Strings(all)
	return all
}

func newRucksack(s string) (rucksack, error) {
	contents := strings.Split(s, "")
	if len(contents)%2 != 0 {
		return rucksack{}, errors.New("invalid rucksack has odd-count of items")
	}
	half := len(contents) / 2
	c1, c2 := contents[0:half], contents[half:]
	if len(c1) != len(c2) {
		return rucksack{}, errors.New("idiot math on slice splitting")
	}
	sort.Strings(c1)
	sort.Strings(c2)
	return rucksack{c1, c2}, nil
}

func (r rucksack) getOverlappingItems() map[string]int {
	overlap := map[string]int{}
	for len(r.c1) > 0 && len(r.c2) > 0 {
		if r.c1[0] != r.c2[0] {
			if r.c1[0] < r.c2[0] {
				if len(r.c1) == 1 {
					break
				}
				r.c1 = r.c1[1:]
			} else {
				if len(r.c2) == 1 {
					break
				}
				r.c2 = r.c2[1:]
			}
			continue
		}
		match := r.c1[0]
		count := 0
		addAndCut := func(c []string) []string {
			s := strings.Join(c, "")
			i := strings.LastIndex(s, match) + 1
			count += i
			if i == len(c) {
				return []string{}
			}
			return c[i:]
		}
		r.c1 = addAndCut(r.c1)
		r.c2 = addAndCut(r.c2)
		overlap[match] = count
	}
	return overlap
}

func intersectOverlaps(o1, o2 map[string]int) map[string]int {
	overlap := map[string]int{}
	for match, count1 := range o1 {
		if count2, ok := o2[match]; ok {
			overlap[match] = count1 + count2
		}
	}
	return overlap
}

func scoreOverlap(o map[string]int) int {
	// _ takes the 0 index so that index is directly score
	values := "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	score := 0
	for letter := range o {
		score += strings.Index(values, letter)
	}
	return score
}

func day3(in io.Reader) (int, int, error) {
	r := bufio.NewReader(in)
	var sacks []rucksack
	for s, err := r.ReadString('\n'); err != io.EOF; s, err = r.ReadString('\n') {
		if err != nil {
			return 0, 0, err
		}
		sack, err := newRucksack(strings.TrimSuffix(s, "\n"))
		if err != nil {
			return 0, 0, err
		}
		sacks = append(sacks, sack)
	}

	sackScore := 0
	for _, sack := range sacks {
		sackScore += scoreOverlap(sack.getOverlappingItems())
	}

	groupScore := 0
	for i := 0; i < len(sacks); i += 3 {
		fauxSack1 := rucksack{sacks[i].allItems(), sacks[i+1].allItems()}
		fauxSack2 := rucksack{sacks[i+1].allItems(), sacks[i+2].allItems()}
		groupScore += scoreOverlap(intersectOverlaps(fauxSack1.getOverlappingItems(), fauxSack2.getOverlappingItems()))
	}
	return sackScore, groupScore, nil
}

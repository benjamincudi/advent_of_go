package a2022

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func anyOverlap(s []string) bool {
	seen := map[string]bool{}
	for _, letter := range s {
		if seen[letter] {
			return true
		}
		seen[letter] = true
	}
	return false
}

type res struct {
	startOfPacket, startOfMessage int
}

type overlapChecker struct {
	duplicates, length int
	state              map[string]int
	history            []string
}

func newOverlapChecker(l int) overlapChecker {
	return overlapChecker{length: l, state: map[string]int{}}
}

func (o *overlapChecker) initFrom(src []string) error {
	if len(src) != o.length {
		return errors.New("initFrom src must match overlapChecker length")
	}
	o.history = src
	for _, s := range src {
		if o.state[s] > 0 {
			o.duplicates++
		}
		o.state[s]++
	}
	return nil
}

func (o *overlapChecker) next(s string) {
	if o.state[s] > 0 {
		o.duplicates++
	}
	o.state[s]++
	if len(o.history) == o.length {
		popped, rest := o.history[0], o.history[1:]
		o.history = append(rest, s)
		o.state[popped]--
		if o.state[popped] > 0 {
			o.duplicates--
		}
	} else if len(o.history) < o.length {
		o.history = append(o.history, s)
	}
}

func (o *overlapChecker) done() bool {
	return len(o.history) == o.length && o.duplicates == 0
}

func day6(in io.Reader) res {
	r := bufio.NewScanner(in)
	r.Scan()
	d := strings.Split(r.Text(), "")
	sop, som := -1, -1
	for i := 4; i <= len(d); i++ {
		if !anyOverlap(d[i-4 : i]) {
			sop = i
			break
		}
	}
	// if we didn't find an answer, a longer answer can't exist
	if sop == -1 {
		return res{-1, -1}
	}
	for i := maxInt(sop+10, 14); i <= len(d); i++ {
		if !anyOverlap(d[i-14 : i]) {
			som = i
			break
		}
	}
	return res{
		startOfPacket:  sop,
		startOfMessage: som,
	}
}

func day6fast(in io.Reader) res {
	r := bufio.NewScanner(in)
	r.Scan()
	d := strings.Split(r.Text(), "")
	sop, som := -1, -1
	sopChecker, somChecker := newOverlapChecker(4), newOverlapChecker(14)
	if err := sopChecker.initFrom(d[0:4]); err != nil {
		panic(err)
	}
	for i := 4; i < len(d); i++ {
		sopChecker.next(d[i])
		if sopChecker.done() {
			sop = i
			break
		}
	}
	// if we didn't find an answer, a longer answer can't exist
	if sop == -1 {
		return res{-1, -1}
	}
	somStartingI := maxInt(sop+10, 14)
	if err := somChecker.initFrom(d[somStartingI-14 : somStartingI]); err != nil {
		panic(err)
	}
	for i := somStartingI; i < len(d); i++ {
		somChecker.next(d[i])
		if somChecker.done() {
			som = i
			break
		}
	}
	return res{
		// answers are 1-indexed
		startOfPacket:  sop + 1,
		startOfMessage: som + 1,
	}
}

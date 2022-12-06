package a2022

import (
	"bufio"
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

func day6(in io.Reader) res {
	r := bufio.NewScanner(in)
	r.Scan()
	d := strings.Split(r.Text(), "")
	sop, som := -1, -1
	for i := 4; i < len(d); i++ {
		if !anyOverlap(d[i-4 : i]) {
			sop = i
			break
		}
	}
	for i := maxInt(sop+10, 14); i < len(d); i++ {
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

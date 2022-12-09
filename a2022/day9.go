package a2022

import (
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
)

type coordinates struct {
	x, y int
}

type ropeInstruction struct {
	D     string
	Count int
}

type ropeKnot struct {
	position coordinates
	history  map[coordinates]bool
	tail     *ropeKnot
}

func makeRopeKnot(i, count int) *ropeKnot {
	if i == count {
		return nil
	}
	return &ropeKnot{coordinates{0, 0}, map[coordinates]bool{coordinates{0, 0}: true}, makeRopeKnot(i+1, count)}
}

func (r *ropeKnot) move(d string) {
	switch d {
	case "R":
		r.position = coordinates{r.position.x + 1, r.position.y}
	case "L":
		r.position = coordinates{r.position.x - 1, r.position.y}
	case "U":
		r.position = coordinates{r.position.x, r.position.y + 1}
	case "D":
		r.position = coordinates{r.position.x, r.position.y - 1}
	}
	r.history[r.position] = true
	r.tail.follow(r.position)
}

func (r *ropeKnot) follow(p coordinates) {
	deltaX, deltaY := p.x-r.position.x, p.y-r.position.y
	dX, dY := abs(deltaX), abs(deltaY)
	if dX > 2 || dY > 2 {
		panic("delta should never exceed 2")
	}
	if dX == 2 || dY == 2 {
		moveX, moveY := sign(deltaX), sign(deltaY)
		r.position = coordinates{r.position.x + moveX, r.position.y + moveY}
		r.history[r.position] = true
	}
	if r.tail != nil {
		r.tail.follow(r.position)
	}
}

func (r *ropeKnot) getTailHistoryLength() int {
	if r.tail != nil {
		return r.tail.getTailHistoryLength()
	}
	return len(r.history)
}

func day9(in io.Reader) (int, int) {
	r := csv.NewReader(in)
	r.Comma = ' '
	var ropeSteps []ropeInstruction
	if err := gocsv.UnmarshalCSVWithoutHeaders(r, &ropeSteps); err != nil {
		panic(err)
	}

	head := makeRopeKnot(0, 2)
	h := makeRopeKnot(0, 10)
	for _, s := range ropeSteps {
		for c := 0; c < s.Count; c++ {
			head.move(s.D)
			h.move(s.D)
		}
	}

	return head.getTailHistoryLength(), h.getTailHistoryLength()
}

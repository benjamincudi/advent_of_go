package a2022

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/gocarina/gocsv"
)

type rpsPlay int

const (
	rpsRock rpsPlay = iota
	rpsPaper
	rpsScissor
)
const (
	lose int = 0
	draw int = 3
	win  int = 6
)

type rpsHand struct {
	play rpsPlay
}

func (r rpsPlay) score() int {
	return int(r) + 1
}

func (r *rpsHand) UnmarshalCSV(s string) error {
	switch s {
	case "A", "X":
		r.play = rpsRock
	case "B", "Y":
		r.play = rpsPaper
	case "C", "Z":
		r.play = rpsScissor
	default:
		return fmt.Errorf("unknown input: %s", s)
	}
	return nil
}

type rpsRound struct {
	Opp, Me rpsHand
}

func (r rpsRound) result(isPlay bool) int {
	keyWins := map[rpsHand]rpsPlay{
		{rpsRock}:    rpsScissor,
		{rpsScissor}: rpsPaper,
		{rpsPaper}:   rpsRock,
	}
	keyLoses := map[rpsHand]rpsPlay{
		{rpsScissor}: rpsRock,
		{rpsPaper}:   rpsScissor,
		{rpsRock}:    rpsPaper,
	}
	if isPlay {
		switch true {
		case r.Opp == r.Me:
			return draw + r.Me.play.score()
		case keyWins[r.Me] == r.Opp.play:
			return win + r.Me.play.score()
		default:
			return lose + r.Me.play.score()
		}
	}
	switch r.Me.play {
	case rpsRock: // Rock == lose
		return lose + keyWins[r.Opp].score()
	case rpsPaper: // Paper == draw
		return draw + r.Opp.play.score()
	default: // Scissor = win
		return win + keyLoses[r.Opp].score()
	}
}
func day2(in io.Reader, isPlay bool) int {
	r := csv.NewReader(in)
	r.Comma = ' '
	var gameRounds []rpsRound
	if err := gocsv.UnmarshalCSVWithoutHeaders(r, &gameRounds); err != nil {
		panic(err)
	}

	total := 0
	for _, round := range gameRounds {
		score := round.result(isPlay)
		total += score
	}
	return total
}

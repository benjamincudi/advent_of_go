package a2022

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
)

type RPSPlay int

const (
	RPSRock RPSPlay = iota
	RPSPaper
	RPSScissor
)
const (
	lose int = 0
	draw     = 3
	win      = 6
)

type RPSHand struct {
	play RPSPlay
}

func (r *RPSHand) Score() int {
	return int(r.play) + 1
}

func (r *RPSHand) UnmarshalCSV(s string) error {
	switch s {
	case "A", "X":
		r.play = RPSRock
	case "B", "Y":
		r.play = RPSPaper
	case "C", "Z":
		r.play = RPSScissor
	default:
		return errors.New(fmt.Sprintf("Unknown input: %s", s))
	}
	return nil
}

type RPSRound struct {
	Opp RPSHand `csv:"lmao"`
	Me  RPSHand `cav:"notthere"`
}

func (r RPSRound) Result(isPlay bool) int {
	keyBeats := map[RPSHand]RPSHand{
		{RPSRock}:    {RPSScissor},
		{RPSScissor}: {RPSPaper},
		{RPSPaper}:   {RPSRock},
	}
	keyLoses := map[RPSHand]RPSHand{
		{RPSScissor}: {RPSRock},
		{RPSPaper}:   {RPSScissor},
		{RPSRock}:    {RPSPaper},
	}
	if isPlay {
		switch true {
		case r.Opp == r.Me:
			return draw + r.Me.Score()
		case keyBeats[r.Me] == r.Opp:
			return win + r.Me.Score()
		default:
			return lose + r.Me.Score()
		}
	}
	switch r.Me.play {
	case RPSRock: // Rock == lose
		shouldPlay := keyBeats[r.Opp]
		return lose + shouldPlay.Score()
	case RPSPaper: // Paper == draw
		return draw + r.Opp.Score()
	default: // Scissor = win
		shouldPlay := keyLoses[r.Opp]
		return win + shouldPlay.Score()
	}
}
func day2(in io.Reader, isPlay bool) int {
	r := csv.NewReader(in)
	r.Comma = ' '
	var gameRounds []RPSRound
	if err := gocsv.UnmarshalCSVWithoutHeaders(r, &gameRounds); err != nil {
		panic(err)
	}

	total := 0
	for _, r := range gameRounds {
		score := r.Result(isPlay)
		total += score
	}
	return total
}

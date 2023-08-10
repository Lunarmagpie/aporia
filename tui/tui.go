package tui

import (
	"reflect"

	"golang.org/x/crypto/ssh/terminal"
)

type Tui struct {
	TermSize TermSize
	Username string
	Password string
	position int
	fields   []field
}

type TermSize struct {
	Lines int
	Cols  int
}

func New() (Tui, error) {
	cols, lines, err := terminal.GetSize(0)

	if err != nil {
		return Tui{}, err
	}

	return Tui{
		TermSize: TermSize{
			Lines: lines,
			Cols:  cols,
		},
		position: 0,
		fields: []field{
			newInput("username"),
			newInput("password"),
			newButton("login"),
		},
	}, nil
}

func (tui *Tui) NextPosition() {
	tui.position = minInt(tui.position+1, len(tui.fields)-1)
}

func (tui *Tui) PrevPosition() {
	tui.position = maxInt(tui.position-1, 0)
}

func (tui *Tui) HandleInput(symbol []int) {
	// Up arrow
	if reflect.DeepEqual(symbol, []int{27, 91, 65}) {
		tui.PrevPosition()
		return
	}
	// Down arrow
	if reflect.DeepEqual(symbol, []int{27, 91, 66}) {
		tui.NextPosition()
		return
	}

	tui.fields[tui.position].onInput(tui, symbol)
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func minInt(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

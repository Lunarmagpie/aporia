package tui

import (
	"fmt"
	"os"
	"reflect"

	"aporia/pam"

	"golang.org/x/term"
)

type Tui struct {
	TermSize TermSize
	position int
	message  string
	fields   []field
	asciiContext asciiArt
	shouldBeCleared bool
}

type TermSize struct {
	Lines int
	Cols  int
}

// Create a new UI. Clears the terminal.
func New() (Tui, error) {
	cols, lines, err := term.GetSize(0)

	if err != nil {
		return Tui{}, err
	}

	fmt.Print("\033[3J")

	return Tui{
		TermSize: TermSize{
			Lines: lines,
			Cols:  cols,
		},
		position: 0,
		message:  "SATA ANDAGI",
		fields: []field{
			newInput("username", false),
			newInput("password", true),
		},
	}, nil
}

func (self *Tui) reset() {
	// TODO
	fmt.Print("Should reset here")
}

func (self *Tui) NextPosition() {
	self.position = minInt(self.position+1, len(self.fields)-1)
}

func (self *Tui) PrevPosition() {
	self.position = maxInt(self.position-1, 0)
}

func (self *Tui) onLastPosition() bool {
	return self.position == len(self.fields)-1
}

func (self *Tui) HandleInput(symbol []int) {
	// Up arrow
	if reflect.DeepEqual(symbol, []int{27, 91, 65}) {
		self.PrevPosition()
		return
	}
	// Down arrowself
	if reflect.DeepEqual(symbol, []int{27, 91, 66}) {
		self.NextPosition()
		return
	}

	// Enter key
	if reflect.DeepEqual(symbol, []int{13}) {
		if self.onLastPosition() {
			self.login()
		} else {
			self.NextPosition()
		}
		return
	}

	// Control + C
	if reflect.DeepEqual(symbol, []int{3}) {
		os.Exit(1)
	}

	self.fields[self.position].onInput(self, symbol)
}

func (self *Tui) login() {
	username := self.fields[0].getContents()
	password := self.fields[1].getContents()

	err := pam.Authenticate(username, password)

	if err != nil {
		self.message = fmt.Sprint(err)
	} else {
		self.message = "Success!"
	}

	self.reset()
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

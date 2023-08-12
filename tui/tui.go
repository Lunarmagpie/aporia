package tui

import (
	"fmt"
	"os"
	"reflect"

	"aporia/pam"

	"golang.org/x/term"
)

type Tui struct {
	TermSize         TermSize
	position         int
	message          string
	fields           []field
	asciiContext     asciiArt
	shouldBeRedrawn  bool
	lastDrawnMessage string
	oldState         *term.State
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

	self := Tui{
		TermSize: TermSize{
			Lines: lines,
			Cols:  cols,
		},
		position:        0,
		message:         "SATA ANDAGI",
		fields:          getFields(),
		shouldBeRedrawn: true,
	}

	self.setup()

	return self, nil
}

// Create the list of fields
func getFields() []field {
	return []field{
		newInput("username", false),
		newInput("password", true),
	}
}

// Functions that need to be called to get the terminal into
// The proper state.
func (self *Tui) setup() {
	self.position = 0
	self.oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
}

func (self *Tui) reset() {
	self.shouldBeRedrawn = true
	self.position = 0
	self.fields = getFields()
	self.setup()
	self.Draw()
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
	// On login, we have to clear the terminal.
	self.shouldBeRedrawn = true

	username := self.fields[0].getContents()
	password := self.fields[1].getContents()

	term.Restore(int(os.Stdin.Fd()), self.oldState)
	err := pam.Authenticate(username, password)

	if err != nil {
		self.message = fmt.Sprint(err)
	} else {
		fmt.Print("\033[H\033[0J")
		self.reset()
		self.message = "Success!"
	}
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

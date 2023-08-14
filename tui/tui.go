package tui

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"

	"aporia/login"

	"golang.org/x/term"
)

type Tui struct {
	TermSize         TermSize
	position         int
	message          string
	fields           []field
	asciiArt         AsciiArt
	shouldBeRedrawn  bool
	lastDrawnMessage string
	loggedIn         bool
	oldState         *term.State
	sessions         []login.Session
}

type TermSize struct {
	Lines int
	Cols  int
}

// Create a new UI. Clears the terminal.
func New(sessions []login.Session) (Tui, error) {
	cols, lines, err := term.GetSize(0)

	if err != nil {
		return Tui{}, err
	}

	state, err := term.GetState(int(os.Stdin.Fd()))

	if err != nil {
		return Tui{}, err
	}

	return Tui{
		TermSize: TermSize{
			Lines: lines,
			Cols:  cols,
		},
		position:        0,
		message:         "SATA ANDAGI",
		fields:          getFields(sessions),
		shouldBeRedrawn: true,
		loggedIn:        false,
		oldState:        state,
		sessions:        sessions,
	}, nil
}

func (self *Tui) Start(charReader CharReader) {
	self.reset()
	self.draw()

	for {
		symbol, err := charReader()

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		self.handleInput(symbol)
		if self.loggedIn {
			break
		}
		self.draw()
	}
}

func (self *Tui) SetAsciiArt(asciiArt AsciiArt) {
	self.asciiArt = asciiArt
	self.message = asciiArt.messages[rand.Intn(len(asciiArt.messages))]
}

// Create the list of fields
func getFields(sessions []login.Session) []field {
	sessionNames := []string{}
	for _, session := range sessions {
		sessionNames = append(sessionNames, session.Name)
	}

	return []field{
		newPicker(sessionNames),
		newInput("username", false),
		newInput("password", true),
	}
}

// Functions that need to be called to get the terminal into
// The proper state.
func (self *Tui) reset() {
	self.loggedIn = false
	self.position = 0
	term.MakeRaw(int(os.Stdin.Fd()))
	self.fields = getFields(self.sessions)
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

func (self *Tui) handleInput(symbol []int) {
	// Up arrow
	if reflect.DeepEqual(symbol, []int{27, 91, 65}) {
		self.PrevPosition()
		return
	}
	// Down arrow and tab
	if reflect.DeepEqual(symbol, []int{27, 91, 66}) || reflect.DeepEqual(symbol, []int{9}) {
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

	sessionName := self.fields[0].getContents()
	username := self.fields[1].getContents()
	password := self.fields[2].getContents()

	term.Restore(int(os.Stdin.Fd()), self.oldState)

	var session login.Session
	for _, this_session := range self.sessions {
		if this_session.Name == sessionName {
			session = this_session
			break
		}
	}

	err := login.Authenticate(username, password, session)

	// We reset the terminal no matter if the login was right or wrong.
	// This way wrong logins make the user re-enter the username and password.
	self.reset()

	if err != nil {
		self.message = fmt.Sprint(err)
	} else {
		self.message = "Success!"
		self.loggedIn = true
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

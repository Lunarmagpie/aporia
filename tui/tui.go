package tui

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"

	"aporia/ansi"
	"aporia/config"
	"aporia/constants"
	"aporia/login"

	"golang.org/x/term"
)

type Tui struct {
	config           config.Config
	TermSize         TermSize
	position         int
	message          string
	fields           []field
	asciiArt         config.AsciiArt
	lastDrawnMessage string
	loggedIn         bool
	termState        term.State
}

type TermSize struct {
	Lines int
	Cols  int
}

// Create a new UI. Clears the terminal.
func New(config config.Config, resetTo term.State) (Tui, error) {
	term.Restore(int(os.Stdin.Fd()), &resetTo)
	ansi.Clear()
	cols, lines, err := term.GetSize(0)

	if err != nil {
		return Tui{}, err
	}

	if err != nil {
		return Tui{}, err
	}

	self := Tui{
		TermSize: TermSize{
			Lines: lines,
			Cols:  cols,
		},
		position:  0,
		message:   "This value should be set before being read.",
		loggedIn:  false,
		termState: resetTo,
		config:    config,
	}
	self.fields = self.getFields()
	return self, nil
}

func (self *Tui) Start() {
	self.setupDraw()
	self.draw()
	_, _ = term.MakeRaw(int(os.Stdin.Fd()))
	charReader := ReadTermChars()

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

func (self *Tui) SetAsciiArt(asciiArt config.AsciiArt) {
	self.asciiArt = asciiArt
	self.message = asciiArt.Messages[rand.Intn(len(asciiArt.Messages))]
}

// Create the list of fields
func (self *Tui) getFields() []field {
	sessionNames := []string{}
	lastSessionIndex := -1

	for i, session := range self.config.Sessions {
		sessionNames = append(sessionNames, session.Name)
		if self.config.LastSession != nil && self.config.LastSession.SessionName == session.Name {
			lastSessionIndex = i
		}
	}

	wmInput := newPicker(sessionNames)
	userInput := newInput("username", false)
	passwdInput := newInput("password", true)

	if lastSessionIndex > -1 {
		wmInput.selected = lastSessionIndex
		self.position = 1
	}
	if self.config.LastSession != nil {
		userInput.contents = self.config.LastSession.User
		self.position = 2
	}

	return []field{
		wmInput,
		userInput,
		passwdInput,
	}
}

func (self *Tui) failedPasswordReset() {
	self.fields[2] = newInput("password", true)
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

	// F11
	if reflect.DeepEqual(symbol, []int{27, 91, 50, 51, 126}) {
		exec.Command(constants.ShutdownCommand[0], constants.ShutdownCommand[1:]...).Run()
		return
	}

	// F12
	if reflect.DeepEqual(symbol, []int{27, 91, 50, 52, 126}) {
		exec.Command(constants.RebootCommand[0], constants.RebootCommand[1:]...).Run()
		return
	}

	// Control + C
	if reflect.DeepEqual(symbol, []int{3}) {
		os.Exit(1)
	}

	self.fields[self.position].onInput(self, symbol)
}

func (self *Tui) login() {
	sessionName := self.fields[0].getContents()
	username := self.fields[1].getContents()
	password := self.fields[2].getContents()

	var session config.Session
	for _, this_session := range self.config.Sessions {
		if this_session.Name == sessionName {
			session = this_session
			break
		}
	}

	self.message = "Authenticating..."
	self.draw()
	term.Restore(int(os.Stdin.Fd()), &self.termState)
	err := login.Authenticate(username, password, session)

	if err != nil {
		// Reset the fields when the password is wrong.
		self.failedPasswordReset()
		self.message = fmt.Sprint(err)
		_, _ = term.MakeRaw(int(os.Stdin.Fd()))
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

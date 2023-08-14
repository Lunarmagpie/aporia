package tui

import (
	"fmt"
	"reflect"
	"strings"
)

type field interface {
	draw(boxSize int) (string, int)
	onInput(tui *Tui, symbol []int)
	getContents() string
}

type input struct {
	prompt   string
	contents string
	masked   bool
}

func (self *input) draw(boxSize int) (string, int) {
	var contents string
	if self.masked {
		contents = strings.Repeat("*", len(self.contents))
	} else {
		contents = self.contents
	}
	return self.prompt + ": " + contents, len(self.prompt) + len(contents) + 3
}

func (self *input) onInput(tui *Tui, symbol []int) {
	if len(symbol) == 1 {
		// Backspace
		if symbol[0] == 127 {
			if len(self.contents) > 0 {
				self.contents = self.contents[:len(self.contents)-1]
			}
			// Other characters
		} else {
			self.contents = self.contents + string(rune(symbol[0]))
		}
	}
}

func (self *input) getContents() string {
	return self.contents
}

func newInput(prompt string, masked bool) *input {
	return &input{
		prompt:   prompt,
		contents: "",
		masked:   masked,
	}
}

// Picker for WMs
// looks like `< name     >`
type picker struct {
	options  []string
	selected int
}

func (self *picker) draw(boxSize int) (string, int) {
	sessionName := self.options[self.selected]

	leftover := boxSize - 6
	afterName := leftover - len(sessionName)

	return fmt.Sprint(" < ", sessionName, strings.Repeat(" ", afterName), " > "), 4
}

func (self *picker) onInput(tui *Tui, symbol []int) {
	// Right arrow
	if reflect.DeepEqual(symbol, []int{27, 91, 67}) {
		self.selected += 1
		if self.selected >= len(self.options) {
			self.selected = 0
		}
		// Left arrow
	} else if reflect.DeepEqual(symbol, []int{27, 91, 68}) {
		self.selected -= 1
		if self.selected < 0 {
			self.selected = len(self.options) - 1
		}
	}
}

func (self *picker) getContents() string {
	return self.options[self.selected]
}

// Options should have a length of at least one.
func newPicker(options []string) *picker {
	return &picker{
		options:  options,
		selected: 0,
	}
}

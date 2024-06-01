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
	location int
	masked   bool
}

func (self *input) draw(boxSize int) (string, int) {
	var contents string
	if self.masked {
		contents = strings.Repeat("*", len(self.contents))
	} else {
		contents = self.contents
	}
	return self.prompt + ": " + contents, len(self.prompt) + self.location + 3
}

func (self *input) onInput(tui *Tui, symbol []int) {
	// Backspace
	if symbol[0] == 127 {
		if self.location > 0 {
			self.contents = self.contents[:self.location-1] + self.contents[self.location:]
			self.location -= 1
		}
		// Right Arrow
	} else if reflect.DeepEqual(symbol, []int{27, 91, 67}) {
		self.location += 1
		if self.location >= len(self.contents) {
			self.location = len(self.contents)
		}
		// Left Arrow
	} else if reflect.DeepEqual(symbol, []int{27, 91, 68}) {
		self.location -= 1
		if self.location < 0 {
			self.location = 0
		}
	} else if len(symbol) == 1 {
		// Other characters
		self.location += 1
		self.contents = self.contents[:self.location-1] + string(rune(symbol[0])) + self.contents[self.location-1:]
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
// looks like `<     name     >`
type picker struct {
	options  []string
	selected int
}

func (self *picker) draw(boxSize int) (string, int) {
	sessionName := self.options[self.selected]

	sideSpacing := (boxSize / 2 - 15)
	leftover := boxSize - 4 - (boxSize / 2 - 15)

	spacingBefore := (leftover - len(sessionName)) / 2
	spacingAfter := leftover - len(sessionName) - spacingBefore

	return fmt.Sprint(strings.Repeat(" ", sideSpacing), " <", strings.Repeat(" ", spacingBefore), sessionName, strings.Repeat(" ", spacingAfter), "> ", strings.Repeat(" ", sideSpacing)), 2 + sideSpacing
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

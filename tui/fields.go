package tui

import "strings"

type field interface {
	draw() string
	onInput(tui *Tui, symbol []int)
	getContents() string
}

type input struct {
	prompt   string
	contents string
	masked bool
}

func (self *input) draw() string {
	var contents string
	if self.masked {
		contents = strings.Repeat("*", len(self.contents))
	} else {
		contents = self.contents
	}
	return self.prompt + ": " + contents
}

func (self *input) onInput(tui *Tui, symbol []int) {
	if len(symbol) == 1 {
		if symbol[0] == 127 {
			if len(self.contents) > 0 {
				self.contents = self.contents[:len(self.contents)-1]
			}
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
		masked: masked,
	}
}

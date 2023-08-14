package tui

import (
	"aporia/ansi"
	"fmt"
	"strings"
)

const horizontal = "─"
const vertical = "│"

const tlCorner = "┌"
const trCorner = "┐"
const blCorner = "└"
const brCorner = "┘"

const boxHeight = 6
const boxWidth = 30

func (self *Tui) draw() error {
	// Reset cursor position
	if self.shouldBeRedrawn {
		ansi.Clear()
		self.shouldBeRedrawn = false
		ansi.MoveCursor(0, 0)
		self.asciiArt.draw(self.TermSize)
		ansi.MoveCursor(0, 0)
		self.drawBox()
	}

	// Draw the message
	if self.lastDrawnMessage != self.message {
		self.lastDrawnMessage = self.message
		ansi.MoveCursor(2, 0)
		self.drawLine(self.message)
	}

	// Draw the currently selected field
	thisLine, cursorPos := self.fields[self.position].draw(boxWidth - 2)

	ansi.MoveCursor(self.position+3, 0)
	self.drawLine(thisLine)

	ansi.MoveCursor(self.position+3, cursorPos + 1)

	return nil
}

// Draw the vertical margin.
func drawMargin(height int) {
	for i := 0; i < (height-boxHeight)/2; i++ {
		fmt.Print("\n\r")
	}
}

func eraseLine(num int) {
	fmt.Print("\033[", num, "K")
}

// Draw the box. Return the vertical lines taken up.
func (self Tui) drawBox() {
	fmt.Print(tlCorner, strings.Repeat(horizontal, boxWidth-2), trCorner, "\n\r")

	self.drawLine(self.message)

	for _, field := range self.fields {
		line, _ := field.draw(boxWidth - 2)
		self.drawLine(line)
	}

	fmt.Print(blCorner, strings.Repeat(horizontal, boxWidth-2), brCorner, "\n\r")
}

func (self Tui) drawLine(text string) {
	fmt.Print(vertical)
	fmt.Print(text)
	fmt.Print(strings.Repeat(" ", maxInt(boxWidth-2-len(text), 0)))
	fmt.Print(vertical, "\n\r")
}

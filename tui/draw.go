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

func (self Tui) Draw() error {
	// drawMargin(self.TermSize.Lines)
	self.drawBox(self.TermSize.Cols)
	// drawMargin(self.TermSize.Lines)

	return nil
}

// Draw the vertical margin.
func drawMargin(height int) {
	for i := 0; i < (height-boxHeight)/2; i++ {
		fmt.Print("\n\r")
	}
}

func (self Tui) drawBox(width int) {
	fmt.Print(tlCorner, strings.Repeat(horizontal, boxWidth-2), trCorner, "\n\r")

	self.drawLine(self.message, width, false)

	for i, field := range self.fields {
		self.drawLine(field.draw(), width, i == self.position)
	}

	fmt.Print(blCorner, strings.Repeat(horizontal, boxWidth-2), brCorner, "\n\r")
}

func (self Tui) drawLine(text string, width int, underline bool) {
	fmt.Print(vertical)
	if underline {
		ansi.Underline()
	}
	fmt.Print(text)
	fmt.Print(strings.Repeat(" ", maxInt(boxWidth-2-len(text), 0)))
	ansi.Reset()
	fmt.Print(vertical, "\n\r")
}

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

// Draw the background, which should only be drawn once.
func (self *Tui) setupDraw(boxWidth int) {
	ansi.Clear()
	ansi.MoveCursor(0, 0)
	draw(self.asciiArt, self.TermSize)
	ansi.MoveCursor(0, 0)
	self.drawTopLeftBox(boxWidth)
	self.drawBox(0, self.TermSize.Cols-13, 14, []string{"F11 Shutdown", "F12 Reboot"})
}

func (self *Tui) draw(boxWidth int) error {
	// Draw the message
	if self.lastDrawnMessage != self.message {
		self.lastDrawnMessage = self.message
		self.drawLine(self.message, 2, 0, boxWidth)
	}

	// Draw the currently selected field
	thisLine, cursorPos := self.fields[self.position].draw(boxWidth - 2)

	self.drawLine(thisLine, self.position+3, 0, boxWidth)

	ansi.MoveCursor(self.position+3, cursorPos+1)

	return nil
}

func eraseLine(num int) {
	fmt.Print("\033[", num, "K")
}

// Draw the box. Return the vertical lines taken up.
func (self Tui) drawTopLeftBox(boxWidth int) {
	lines := []string{""}
	for _, field := range self.fields {
		line, _ := field.draw(boxWidth - 2)
		lines = append(lines, line)
	}

	self.drawBox(0, 0, boxWidth, lines)

}

func (self Tui) drawBox(line int, col int, width int, text []string) {
	ansi.MoveCursor(line, col)
	fmt.Print(tlCorner, strings.Repeat(horizontal, width-2), trCorner)

	for i, line := range text {
		self.drawLine(line, i+2, col, width)
	}

	ansi.MoveCursor(len(text)+2+line, col)
	fmt.Print(blCorner, strings.Repeat(horizontal, width-2), brCorner)
}
func (self Tui) drawLine(text string, line, col, width int) {
	ansi.MoveCursor(line, col)
	fmt.Print(vertical)
	fmt.Print(text)
	fmt.Print(strings.Repeat(" ", maxInt(width-2-len(text), 0)))
	fmt.Print(vertical)
}

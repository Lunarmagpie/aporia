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

func (tui Tui) Draw() error {
	ansi.Clear()

	drawMargin(tui.TermSize.Lines)
	tui.drawBox(tui.TermSize.Cols)
	drawMargin(tui.TermSize.Lines)

	return nil
}

// Draw the vertical margin.
func drawMargin(height int) {
	for i := 0; i < (height-boxHeight)/2; i++ {
		fmt.Println()
	}
}

func (tui Tui) drawBox(width int) {
	fmt.Println(tlCorner + strings.Repeat(horizontal, boxWidth-2) + trCorner)

	for i, field := range tui.fields {
		if i == tui.position {
			ansi.Underline()
		}
		fmt.Printf(field.draw() + "\n")
		ansi.Reset()
	}

	fmt.Println(blCorner + strings.Repeat(horizontal, boxWidth-2) + brCorner)
}

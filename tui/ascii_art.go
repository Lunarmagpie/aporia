package tui

import (
	"aporia/config"
	"fmt"
	"strings"
)

/*
A struct that makes loading and printing ascii art easy.
*/

func draw(self config.AsciiArt, termSize TermSize) {
	linesSkip := maxInt((self.Lines-termSize.Lines)/2, 0)
	colsSkip := maxInt((self.Cols-termSize.Cols)/2, 0)

	startLine := maxInt((termSize.Lines-self.Lines)/2, 0)
	startCol := maxInt((termSize.Cols-self.Cols)/2, 0)

	for i := 0; i < startLine; i++ {
		fmt.Print("\n\r")
	}

	for i := 0; i < termSize.Lines && i+linesSkip < self.Lines; i++ {
		if colsSkip >= len(self.StrLines[i]) {
			continue
		}
		if i != 0 {
			fmt.Print("\n\r")
		}
		maxSize := minInt(termSize.Cols+colsSkip, len(self.StrLines[i+linesSkip]))
		fmt.Print(strings.Repeat(" ", startCol), self.StrLines[i+linesSkip][colsSkip:maxSize])
	}
}

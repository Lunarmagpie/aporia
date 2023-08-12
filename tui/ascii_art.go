package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

/*
A struct that makes loading and printing ascii art easy.
*/
type asciiArt struct {
	lines  []string
	height int
	width  int
}

func newAscii(art string) asciiArt {
	lines := strings.Split(art, "\n")

	longestLine := utf8.RuneCountInString(lines[0])

	for _, line := range lines[1:] {
		if len(line) > longestLine {
			longestLine = utf8.RuneCountInString(line)
		}
	}

	return asciiArt{
		lines:  lines,
		width:  longestLine,
		height: len(lines),
	}
}

func (self *asciiArt) calculatePosition(screenWidth int, screenHeight int) (int, int) {
	x := (screenWidth - self.width) / 2
	y := (screenHeight - self.height) / 2
	return x, y
}

func (self *asciiArt) draw(xOffset int, yOffset int) {
	for i := 0; i < yOffset; i++ {
		fmt.Print("\n\r")
	}

	for _, line := range self.lines {
		fmt.Print(strings.Repeat(" ", xOffset), line, "\n\r")
	}
}

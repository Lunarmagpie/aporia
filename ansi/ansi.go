package ansi

import "fmt"

func Clear() {
	// fmt.Print("\033[H\033[0J")
}

func MoveCursor(line int, col int) {
	fmt.Print("\033[", line, ";", col, "H")
}

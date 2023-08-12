package ansi

import "fmt"

func Clear() {
	fmt.Print("\033[H\033[0J")
}

// Erase n chars from the start of the line
func EraseChars(num int) {
	fmt.Print("\033[", num, "C")
	fmt.Print("\033[1K" )
	fmt.Print("\033[", num, "D")
}

func MoveCursor(line int, col int) {
	fmt.Print("\033[", line, ";", col, "H")
}

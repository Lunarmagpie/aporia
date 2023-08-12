package ansi

import "fmt"

func Clear() {
	fmt.Print("\033[H\033[0J")
}

// Erase n chars in front of the cursor
func EraseChars(num int) {
	fmt.Print("\033[", num, "C")
	fmt.Print("\033[1K" )
	fmt.Print("\033[", num, "D")
}

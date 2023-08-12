package tui

import "fmt"

// Position the cursor to the right place for the user.
func moveCursor(line int, col int) {
	fmt.Print("\033[", line, ";", col, "H")
}

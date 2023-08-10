package ansi

import "fmt"

func Underline() {
	fmt.Print("\033[4m")
}

func Reset() {
	fmt.Print("\033[0m")
}

func Clear() {
	fmt.Print("\033[2J")
}

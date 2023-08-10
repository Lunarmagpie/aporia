package tui

import (
	"os"
	"os/exec"
	"reflect"
)

// Read characters from the terminal. A sequence such as \027[n will be
// returned as one character.
func ReadTermChars() func() ([]int, error) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	buff := []int{}

	readOneChar := func() (int, error) {
		var nextByte []byte = make([]byte, 1)
		_, err := os.Stdin.Read(nextByte)

		if err != nil {
			return 0, err
		}

		keycode := int(nextByte[0])

		return keycode, nil
	}

	return func() ([]int, error) {
		for {
			nextChar, err := readOneChar()

			if err != nil {
				return []int{}, err
			}

			buff = append(buff, nextChar)

			if reflect.DeepEqual(buff, []int{27}) {
				continue
			}

			if reflect.DeepEqual(buff, []int{27, 91}) {
				continue
			}

			copy := buff
			buff = []int{}
			return copy, nil
		}
	}
}

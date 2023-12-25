package tui

import (
	"os"
	"reflect"
)

type CharReader = func() ([]int, error)

// Read characters from the terminal. A sequence such as \027[n will be
// returned as one character.
// Note that this function will not put the terminal into raw mode.
func ReadTermChars() func() ([]int, error) {
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

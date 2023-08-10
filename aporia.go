package main

import (
	"fmt"

	"aporia/tui"
)

func main() {
	ui, _ := tui.New()
	charReader := tui.ReadTermChars()

	ui.Draw()

	for {
		symbol, err := charReader()

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		ui.HandleInput(symbol)
		ui.Draw()
	}
}

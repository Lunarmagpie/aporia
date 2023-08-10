package main

import (
	"fmt"

	"aporia/pam"
	"aporia/tui"
)

func main() {
	ui, _ := tui.New()
	charReader := tui.ReadTermChars()

	err := pam.Authenticate("lunarmagpie", "test")

	if err != nil {
		fmt.Printf(err.Error())
	}

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

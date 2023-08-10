package main

import (
	"fmt"

	"aporia/tui"
)

func main() {
	fmt.Printf("Entering the aporia display manager.")
	
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

package main

import (
	"fmt"
	"math/rand"
	"time"

	"aporia/tui"
)

func main() {
	rand.Seed(time.Now().Unix())

	config, err := loadConfig()
	if err != nil {
		config_ := defaultConfig()
		config = &config_
	}

	ui, _ := tui.New()
	ui.SetAsciiArt(config.randomAscii())
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

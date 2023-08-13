package main

import (
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

	charReader := tui.ReadTermChars()

	for {
		ui.SetAsciiArt(config.randomAscii())
		ui.Start(charReader)
	}
}

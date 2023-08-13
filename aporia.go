package main

import (
	"math/rand"
	"runtime"
	"time"

	"aporia/tui"
)

func main() {
	rand.Seed(time.Now().Unix())
	runtime.LockOSThread()

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

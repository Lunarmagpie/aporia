package main

import (
	"aporia/config"
	"aporia/tui"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	runtime.LockOSThread()

	configObj, err := config.LoadConfig()
	if err != nil {
		config_ := config.DefaultConfig()
		configObj = &config_
	}

	ui, _ := tui.New(*configObj)

	charReader := tui.ReadTermChars()

	for {
		ui.SetAsciiArt(configObj.RandomAscii())
		ui.Start(charReader)
	}
}

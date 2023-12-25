package main

import (
	"aporia/config"
	"aporia/tui"
	"math/rand"
	"os"
	"runtime"
	"time"

	"golang.org/x/term"
)

func main() {
	rand.New(rand.NewSource(time.Now().Unix()))
	runtime.LockOSThread()

	termState, _ := term.GetState(int(os.Stdin.Fd()))

	for {
		configObj, err := config.LoadConfig()
		if err != nil {
			config_ := config.DefaultConfig()
			configObj = &config_
		}
		ui, _ := tui.New(*configObj, *termState)
		ui.SetAsciiArt(configObj.RandomAscii())
		ui.Start()
	}
}

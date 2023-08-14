package main

import (
	"aporia/constants"
	"aporia/login"
	"aporia/tui"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type config struct {
	asciiArts []tui.AsciiArt
	sessions  []login.Session
}

func parseAsciiFile(filepath string) (*tui.AsciiArt, error) {
	contents, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	contentsStr := string(contents)
	contentsLines := strings.Split(contentsStr, "\n")

	messages := []string{"Login:"}
	origin := tui.Center

	asciiLines := []string{}
	var asciiStartLine int

	for i, line := range contentsLines {
		if strings.HasPrefix(line, "messages:") {
			after, _ := strings.CutPrefix(line, "messages:")
			messages = strings.Split(after, ",")
			for i, message := range messages {
				messages[i] = strings.TrimSpace(message)
			}
		}
		if strings.HasPrefix(line, "origin") {
			origin = tui.Center
		}
		if strings.HasPrefix(line, "-") {
			asciiStartLine = i
		}
	}

	for i := asciiStartLine + 1; i < len(contentsLines); i++ {
		if len(strings.TrimSpace(contentsLines[i])) == 0 {
			continue
		}
		asciiLines = append(asciiLines, contentsLines[i])
	}

	ascii := tui.NewAsciiArt(
		strings.Join(asciiLines, "\n"),
		messages,
		origin,
	)

	return &ascii, nil
}

func loadConfig() (*config, error) {
	entries, err := os.ReadDir(constants.ConfigDir)
	if err != nil {
		return nil, err
	}

	asciiArts := []tui.AsciiArt{}
	sessions := []login.Session{}

	// NOTE: We don't need a mutex here because the runtime is locked to one thread.
	wg := sync.WaitGroup{}

	parseEntry := func(entry fs.DirEntry) {
		defer wg.Done()
		filepath := filepath.Join(constants.ConfigDir, entry.Name())
		if strings.HasSuffix(entry.Name(), "."+constants.AsciiFileExt) {
			asciiFile, err := parseAsciiFile(filepath)
			if err != nil {
				return
			}
			asciiArts = append(asciiArts, *asciiFile)
		} else if strings.HasSuffix(entry.Name(), ".x11") {
			name := strings.TrimSuffix(entry.Name(), ".x11")
			sessions = append(sessions, login.X11Session(name, filepath))
		} else if strings.HasSuffix(entry.Name(), ".wayland") {
			name := strings.TrimSuffix(entry.Name(), ".wayland")
			sessions = append(sessions, login.WaylandSession(name, filepath))
		}
	}

	for _, entry := range entries {
		wg.Add(1)
		go parseEntry(entry)
	}

	wg.Wait()

	return &config{
		asciiArts: asciiArts,
		sessions:  append(sessions, login.ShellSession()),
	}, nil
}

// Default config file to be used if the user's config could not be loaded.
func defaultConfig() config {
	return config{
		asciiArts: []tui.AsciiArt{},
		sessions:  []login.Session{login.ShellSession()},
	}
}

func (self *config) randomAscii() tui.AsciiArt {
	if len(self.asciiArts) == 0 {
		return tui.NewAsciiArt(
			constants.DefaultAsciiArt,
			constants.DefaultMessages(),
			tui.Center,
		)
	}
	return self.asciiArts[rand.Intn(len(self.asciiArts))]
}

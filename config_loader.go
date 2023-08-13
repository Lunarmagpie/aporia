package main

import (
	"aporia/constants"
	"aporia/tui"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	asciiArts []tui.AsciiArt
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

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "."+constants.AsciiFileExt) {
			asciiFile, err := parseAsciiFile(filepath.Join(constants.ConfigDir, entry.Name()))
			if err != nil {
				continue
			}
			asciiArts = append(asciiArts, *asciiFile)
		}
	}

	return &config{
		asciiArts: asciiArts,
	}, nil
}

// Default config file to be used if the user's config could not be loaded.
func defaultConfig() config {
	return config{
		asciiArts: []tui.AsciiArt{},
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

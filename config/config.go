package config

import (
	"aporia/constants"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

type Config struct {
	AsciiArts   []AsciiArt
	Sessions    []Session
	LastSession *LastSession
}

type SessionType string

const (
	ShellSession   = "tty"
	X11Session     = "x11"
	WaylandSession = "wayland"
)

type Session struct {
	Name        string
	SessionType SessionType
	// The filepath to the launcher file for the session, or null if its a shell session.
	Filepath *string
}

func newShellSession() Session {
	return Session{
		Name:        "shell",
		SessionType: ShellSession,
		Filepath:    nil,
	}
}

func newX11Session(name string, startxPath string) Session {
	return Session{
		Name:        name,
		SessionType: X11Session,
		Filepath:    &startxPath,
	}
}

func newWaylandSession(name string, filepath string) Session {
	return Session{
		Name:        name,
		SessionType: WaylandSession,
		Filepath:    &filepath,
	}
}

type LastSession struct {
	SessionName string
	User        string
}

type Origin string

const (
	Center Origin = "center"
)

type AsciiArt struct {
	StrLines []string
	Lines    int
	Cols     int

	Messages []string
	Origin   Origin
}

func newAsciiArt(art string, messages []string, origin Origin) AsciiArt {
	lines := strings.Split(art, "\n")

	longestLine := utf8.RuneCountInString(lines[0])

	for _, line := range lines[1:] {
		if len(line) > longestLine {
			longestLine = utf8.RuneCountInString(line)
		}
	}

	return AsciiArt{
		StrLines: lines,
		Cols:     longestLine,
		Lines:    len(lines),
		Messages: messages,
		Origin:   origin,
	}
}

func parseAsciiFile(filepath string) (*AsciiArt, error) {
	contents, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	contentsStr := string(contents)
	contentsLines := strings.Split(contentsStr, "\n")

	messages := []string{"Login:"}
	origin := Center

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
			origin = Center
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

	ascii := newAsciiArt(
		strings.Join(asciiLines, "\n"),
		messages,
		origin,
	)

	return &ascii, nil
}

// Parses a file saved for an old session
func loadLastSession() (*LastSession, error) {
	file, err := os.ReadFile(constants.LastSessionFile)
	if err != nil {
		return nil, err
	}

	contents := strings.Split(string(file), "\n")
	if len(contents) < 2 {
		return nil, errors.New("Last session file was configured incorrectly")
	}

	return &LastSession{
		SessionName: contents[0],
		User:        contents[1],
	}, nil

}

func SaveSession(sessionName string, user string) {
	contents := fmt.Sprint(sessionName, "\n", user)
	os.WriteFile(constants.LastSessionFile, []byte(contents), 0644)
}

func LoadConfig() (*Config, error) {
	entries, err := os.ReadDir(constants.ConfigDir)
	if err != nil {
		return nil, err
	}

	sessions := []Session{}

	// NOTE: We don't need a mutex here because the runtime is locked to one thread.
	wg := sync.WaitGroup{}

	asciiArts := []AsciiArt{}
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
			sessions = append(sessions, newX11Session(name, filepath))
		} else if strings.HasSuffix(entry.Name(), ".wayland") {
			name := strings.TrimSuffix(entry.Name(), ".wayland")
			sessions = append(sessions, newWaylandSession(name, filepath))
		}
	}

	for _, entry := range entries {
		wg.Add(1)
		go parseEntry(entry)
	}

	wg.Wait()

	session, _ := loadLastSession()

	return &Config{
		AsciiArts:   asciiArts,
		Sessions:    append(sessions, newShellSession()),
		LastSession: session,
	}, nil
}

// Default config file to be used if the user's config could not be loaded.
func DefaultConfig() Config {
	return Config{
		AsciiArts: []AsciiArt{},
		Sessions:  []Session{newShellSession()},
	}
}

func (self *Config) RandomAscii() AsciiArt {
	if len(self.AsciiArts) == 0 {
		return newAsciiArt(
			constants.DefaultAsciiArt,
			constants.DefaultMessages(),
			Center,
		)
	}
	return self.AsciiArts[rand.Intn(len(self.AsciiArts))]
}

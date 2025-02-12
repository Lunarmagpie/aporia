package config

import (
	"aporia/constants"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type Config struct {
	AsciiArts    []AsciiArt
	isAsciiError bool
	Sessions     []Session
	LastSession  *LastSession
	Extra        ConfigFile
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
	// Or there is an exec function to be called.
	Exec *string
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

type AsciiArt struct {
	StrLines []string
	Lines    int
	Cols     int
	Messages []string
	name     string
}

func newAsciiArt(name string, art string, messages []string) AsciiArt {
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
		name:     name,
	}
}

func parseAsciiFile(filename string) (*AsciiArt, error) {
	contents, err := os.ReadFile(filepath.Join(constants.ConfigDir, filename))

	if err != nil {
		return nil, err
	}

	contentsStr := string(contents)
	contentsLines := strings.Split(contentsStr, "\n")

	messages := []string{"Enter Credentials:"}

	asciiLines := []string{}
	asciiStartLine := -1

	for i, line := range contentsLines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "messages" && fields[1] == "=" {
			if len(fields) < 3 {
				return nil, errors.New("Formatting error in file: " + filename)
			}
			messages = strings.Split(fields[2], ",")
			for i, message := range messages {
				messages[i] = strings.TrimSpace(message)
			}
		}
		if strings.HasPrefix(line, "---") && asciiStartLine == -1 {
			asciiStartLine = i
		}
	}

	if asciiStartLine == -1 {
		asciiStartLine = 0
	}

	for i := asciiStartLine + 1; i < len(contentsLines); i++ {
		if len(strings.TrimSpace(contentsLines[i])) == 0 {
			continue
		}
		asciiLines = append(asciiLines, contentsLines[i])
	}

	ascii := newAsciiArt(
		strings.TrimSuffix(filename, "."+constants.AsciiFileExt),
		strings.Join(asciiLines, "\n"),
		messages,
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
		return nil, errors.New("Last session file was configured incorrectly.\nRun the command `# rm /etc/aporia/.last-session` to fix.")
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

type ConfigFile struct {
	Ascii               *string
	BoxWidth            int
	ShutdownCommand     []string
	RebootCommand       []string
	XsessionsPath       string
	WaylandSessionsPath string
}

func parseConfigFile() (ConfigFile, error) {
	config := ConfigFile{
		BoxWidth:            constants.BoxWidth,
		ShutdownCommand:     constants.ShutdownCommand,
		RebootCommand:       constants.RebootCommand,
		XsessionsPath:       constants.X11SessionsDir,
		WaylandSessionsPath: constants.WaylandSessionsDir,
	}

	filepath := filepath.Join(constants.ConfigDir, constants.ConfigFile)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return config, nil
	}

	contentsStr := string(data)
	contentsLines := strings.Split(contentsStr, "\n")

	for _, line := range contentsLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) <= 2 {
			return config, errors.New("Error in user config file. Expected expression after =.")
		}

		if fields[0] == "ascii" {
			config.Ascii = &fields[2]
		}
		if fields[0] == "box_wdith" {
			i, err := strconv.Atoi(fields[2])
			if err != nil {
				return config, errors.New("Received non integer for box_width")
			}
			config.BoxWidth = i
		}
		if fields[0] == "shutdown_cmd" {
			config.ShutdownCommand = fields[:2]
		}
		if fields[0] == "reboot_cmd" {
			config.RebootCommand = fields[:2]
		}
		if fields[0] == "xsessions_path" {
			config.XsessionsPath = fields[2]
		}
		if fields[0] == "wayland_sessions_path" {
			config.WaylandSessionsPath = fields[2]
		}
	}

	return config, nil
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
	asciiArtErrors := []error{}
	parseEntry := func(entry fs.DirEntry) {
		defer wg.Done()
		filepath := filepath.Join(constants.ConfigDir, entry.Name())
		if strings.HasSuffix(entry.Name(), "."+constants.AsciiFileExt) {
			asciiFile, err := parseAsciiFile(entry.Name())
			if err != nil {
				asciiArtErrors = append(asciiArtErrors, err)
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

	configFile, configError := parseConfigFile()
	sessions = append(sessions, desktopCrawl(configFile.XsessionsPath, configFile.WaylandSessionsPath)...)

	if configError != nil {
		asciiArtErrors = append(asciiArtErrors, configError)
	}

	isAsciiError := false
	if len(asciiArtErrors) != 0 {
		isAsciiError = true
		asciiArt := ""
		for _, err := range asciiArtErrors {
			asciiArt = asciiArt + fmt.Sprintln(err)
		}
		asciiArts = []AsciiArt{
			newAsciiArt("Error", asciiArt, []string{"There was an error :("}),
		}
	}

	return &Config{
		isAsciiError: isAsciiError,
		Sessions:     append(sessions, newShellSession()),
		LastSession:  session,
		Extra:        configFile,
		AsciiArts:    asciiArts,
	}, nil
}

// Default config file to be used if the user's config could not be loaded.
func DefaultConfig() Config {
	return Config{
		AsciiArts: []AsciiArt{},
		Sessions:  []Session{newShellSession()},
	}
}

func (self *Config) GetAscii() AsciiArt {
	if self.isAsciiError {
		return self.AsciiArts[0]
	}

	if self.Extra.Ascii != nil {
		for _, file := range self.AsciiArts {
			if file.name == *self.Extra.Ascii {
				return file
			}
		}
		return newAsciiArt(
			"This doesn't matter because it is never read.",
			"ascii art `"+*self.Extra.Ascii+"` not found",
			constants.DefaultMessages(),
		)
	}

	if len(self.AsciiArts) == 0 {
		return newAsciiArt(
			"This doesn't matter because it is never read.",
			constants.DefaultAsciiArt,
			constants.DefaultMessages(),
		)
	}

	return self.AsciiArts[rand.Intn(len(self.AsciiArts))]
}

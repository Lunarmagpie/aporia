package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func desktopCrawl(xsessionsPath string, waylandSessionsPath string) []Session {
	return append(
		desktopCrawlPath(xsessionsPath, X11Session),
		desktopCrawlPath(waylandSessionsPath, WaylandSession)...,
	)
}

func desktopCrawlPath(path string, sessionType SessionType) []Session {
	entries, err := os.ReadDir(path)

	if err != nil {
		return []Session{}
	}

	wg := sync.WaitGroup{}
	desktops := []Session{}

	parseEntry := func(path string) {
		defer wg.Done()

		contents_bytes, err := os.ReadFile(path)
		contents := string(contents_bytes)

		if err != nil {
			return
		}

		var name *string
		var exec *string

		for _, line := range strings.Split(contents, "\n") {
			lineClean := strings.TrimSpace(line)

			parts := strings.Split(lineClean, "=")

			if len(parts) < 2 {
				continue
			}

			lClean := strings.TrimSpace(parts[0])
			rClean := strings.TrimSpace(parts[1])

			if lClean == "Name" {
				name = &rClean
			}
			if lClean == "Exec" {
				exec = &rClean
			}
		}

		if name == nil || exec == nil {
			return
		}

		desktops = append(desktops, Session{
			Name: *name,
			Exec: exec,
			SessionType: sessionType,
		})
	}

	for _, entry := range entries {
		filepath := filepath.Join(path, entry.Name())
		wg.Add(1)
		go parseEntry(filepath)
	}

	wg.Wait()

	return desktops

}

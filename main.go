package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	config, err := NewConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	url, err := getRawUrl()
	if err != nil {
		slog.Error("failed to fetch standings URL from GitHub", "err", err)
		os.Exit(1)
	}

	file, err := getRawFile(url)
	if err != nil {
		slog.Error("failed to download standings file", "url", url, "err", err)
		os.Exit(1)
	}

	standings := parseStandings(file)

	synced, err := SyncStandings(config, standings)
	if err != nil {
		slog.Error("failed to sync standings to database", "err", err)
		os.Exit(1)
	}

	if synced {
		slog.Info("standings synced", "count", len(standings))
	} else {
		slog.Info("standings up to date, skipping sync", "date", standings[0].StandingsDate)
	}
}

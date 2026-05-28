package main

import (
	"strconv"
	"strings"
	"time"
)

type VRSEntry struct {
	Standing      int       `bson:"standing"`
	Points        int       `bson:"points"`
	TeamName      string    `bson:"team_name"`
	Roster        []string  `bson:"roster"`
	StandingsDate string    `bson:"standings_date"`
	SyncedAt      time.Time `bson:"synced_at"`
}

// parseStandings parses the contents of a VRS standings markdown file into a slice of VRSEntry.
// The header line ("### Standings as of 2026_05_04<br />") is used to populate StandingsDate.
// A data row looks like: | 1 | 2081 | Vitality | apEX, flameZ, mezii, ropz, ZywOo | [details](...) |
func parseStandings(content string) []VRSEntry {
	var standingsDate string
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if date, ok := strings.CutPrefix(line, "### Standings as of "); ok {
			standingsDate = strings.TrimSpace(strings.TrimSuffix(date, "<br />"))
			break
		}
	}

	seen := make(map[string]bool)
	var entries []VRSEntry

	for line := range strings.SplitSeq(content, "\n") {
		// data rows start with "| 1 |" etc — skip header, separator, empty lines
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}

		cols := strings.Split(line, "|")
		if len(cols) < 5 {
			continue
		}

		standingStr := strings.TrimSpace(cols[1])
		standing, err := strconv.Atoi(standingStr)
		if err != nil {
			continue // skips header and separator rows
		}

		teamName := strings.TrimSpace(cols[3])
		if seen[teamName] {
			continue // keep only first (best) occurrence
		}
		seen[teamName] = true

		roster := strings.Split(strings.TrimSpace(cols[4]), ", ")

		points, _ := strconv.Atoi(strings.TrimSpace(cols[2]))

		entries = append(entries, VRSEntry{
			Standing:      standing,
			Points:        points,
			TeamName:      teamName,
			Roster:        roster,
			StandingsDate: standingsDate,
		})
	}

	return entries
}

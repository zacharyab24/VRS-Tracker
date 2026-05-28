package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// getRawUrl returns the raw.githubusercontent.com URL for the latest VRS global standings file.
// It assumes the current year is 2026 and does not look for standings in other years.
// Files are named with a yyyy_mm_dd suffix, so sorting alphabetically guarantees the last
// entry is always the most recent snapshot.
func getRawUrl() (string, error) {
	const gh_api_url = "https://api.github.com/repos/ValveSoftware/counter-strike_regional_standings/contents/live/2026"
	resp, err := http.Get(gh_api_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var json []map[string]any
	if err := decoder.Decode(&json); err != nil {
		return "", err
	}
	latest := len(json) - 1
	file := json[latest]["download_url"].(string)

	return file, nil
}

// getRawFile fetches the raw markdown content of a VRS standings file from the given URL.
// The URL is expected to be a raw.githubusercontent.com URL obtained from getRawUrl.
func getRawFile(fileUrl string) (string, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

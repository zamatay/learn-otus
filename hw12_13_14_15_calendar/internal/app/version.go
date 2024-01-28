package app

import (
	"encoding/json"
	"os"
)

var (
	release   = "0"
	buildDate = "0"
	gitHash   = "1"
)

func PrintVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		log.Error("error while decode version info: %v\n", err)
	}
}

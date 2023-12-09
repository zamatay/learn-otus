package app

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("error while decode version info: %v\n", err)
	}
}

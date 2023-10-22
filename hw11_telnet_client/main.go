package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"
)

type Options struct {
	timeout time.Duration
	host    string
	port    string
}

var programOptions Options

func init() {
	flag.DurationVar(&programOptions.timeout, "timeout", 100, "timeout duration")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	slog.Info("Flag", slog.Duration("timeout", programOptions.timeout))
	if len(os.Args) < 3 {
		log.Fatal("hostAndPortInvalid")
	}
	programOptions.host = os.Args[len(os.Args)-2]
	programOptions.port = os.Args[len(os.Args)-1]

	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}

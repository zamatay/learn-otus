package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Options struct {
	timeout time.Duration
	host    string
	port    string
	address string
}

var programOptions Options

func init() {
	flag.DurationVar(&programOptions.timeout, "timeout", 100, "timeout duration")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if len(os.Args) < 3 {
		log.Fatal("hostAndPortInvalid")
	}

	programOptions.host = os.Args[len(os.Args)-2]
	programOptions.port = os.Args[len(os.Args)-1]
	run()
}

func run() {
	programOptions.address = fmt.Sprintf("%s:%s", programOptions.host, programOptions.port)
	c := NewTelnetClient(programOptions.address, programOptions.timeout, io.NopCloser(os.Stdin), os.Stdout)
	c.Connect()
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

	defer func() {
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ch:
				return
			default:
				if err := c.Send(); err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ch:
				return
			default:
				if err := c.Receive(); err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}()

	wg.Wait()
}

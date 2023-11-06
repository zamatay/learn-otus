package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"os/signal"
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
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		cancel()
	}()

	g.Go(func() error {
		<-gCtx.Done()
		c.Close()
		return nil
	})

	g.Go(func() error {
		for {
			if err := c.Send(); err != nil {
				log.Print(err)
				break
			}
		}
		return nil
	})

	g.Go(func() error {
		for {
			if err := c.Receive(); err != nil {
				log.Print(err)
				break
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

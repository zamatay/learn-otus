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
		log.Println(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	defer func() {
		log.Println("Close connection")
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				return nil
			default:
				if err := c.Send(); err != nil {
					log.Fatal(err)
					return nil
				}
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				return nil
			default:
				if err := c.Receive(); err != nil {
					log.Fatal(err)
					return nil
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

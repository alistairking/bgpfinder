package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/alistairking/bgpfinder"
)

type ProjectCmd struct {
	// TODO
}

func (p *ProjectCmd) Run() error {
	projs := bgpfinder.Projects()
	for _, proj := range projs {
		fmt.Println(proj)
	}
	return nil
}

type BgpfCLI struct {
	// sub commands
	Project ProjectCmd `cmd help:"Get information about supported projects"`

	// TODO

	// logging configuration
	bgpfinder.LoggerConfig
}

func handleSignals(ctx context.Context, log bgpfinder.Logger, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigCh:
				log.Info().Msgf("Signal recevied, triggering shutdown")
				cancel()
				return
			}
		}
	}()
}

func main() {
	// Parse command line args
	var cliCfg BgpfCLI
	k := kong.Parse(&cliCfg)
	k.Validate()

	// Set up context, logger, and signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logp, err := bgpfinder.NewLogger(cliCfg.LoggerConfig)
	k.FatalIfErrorf(err)
	handleSignals(ctx, *logp, cancel)

	// calls the appropriate command "Run" method
	// TODO: pass some state here (logging?)
	err = k.Run()
	k.FatalIfErrorf(err)

	// Wait a moment for the logger to drain any remaining messages
	time.Sleep(time.Second)
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/alistairking/bgpfinder"
)

type BgpfCLI struct {
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

	// TODO
	logp.Info().Msgf("Here is where the things would happen")

	// Wait a moment for the logger to drain any remaining messages
	time.Sleep(time.Second)
}

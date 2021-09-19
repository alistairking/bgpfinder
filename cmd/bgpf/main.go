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

// TODO: think about how these projects/collectors queries should work.
//
// Other than for interactive exploration, I guess people will want to be able
// to use the output from these commands to drive shell scripts. E.g., list all
// RV collectors and then do something for each one. Or, list all supported
// collectors.

type ProjectsCmd struct {
	// TODO
}

func (p *ProjectsCmd) Run(log bgpfinder.Logger) error {
	projs, err := bgpfinder.Projects()
	if err != nil {
		return fmt.Errorf("Failed to get project list: %v", err)
	}
	for _, proj := range projs {
		fmt.Println(proj)
	}
	return nil
}

type CollectorsCmd struct {
	Project string `help:"Show collectors for the given project"`
}

func (p *CollectorsCmd) Run(log bgpfinder.Logger) error {
	colls, err := bgpfinder.Collectors(p.Project)
	if err != nil {
		return fmt.Errorf("Failed to get collector list: %v", err)
	}
	for _, coll := range colls {
		fmt.Println(coll.String())
	}
	return nil
}

type BgpfCLI struct {
	// sub commands
	Projects   ProjectsCmd   `cmd help:"Get information about supported projects"`
	Collectors CollectorsCmd `cmd help:"Get information about supported collectors"`

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
	err = k.Run(*logp)
	k.FatalIfErrorf(err)

	// Wait a moment for the logger to drain any remaining messages
	time.Sleep(time.Second)
}

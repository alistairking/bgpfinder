package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func (p *ProjectsCmd) Run(log bgpfinder.Logger, cli BgpfCLI) error {
	projs, err := bgpfinder.Projects()
	if err != nil {
		return fmt.Errorf("failed to get project list: %v", err)
	}
	for _, proj := range projs {
		fmt.Println(proj)
	}
	return nil
}

type CollectorsCmd struct {
	Project string `help:"Show collectors for the given project"`
}

func (p *CollectorsCmd) Run(log bgpfinder.Logger, cli BgpfCLI) error {
	colls, err := bgpfinder.Collectors(p.Project)
	if err != nil {
		return fmt.Errorf("failed to get collector list: %v", err)
	}
	for _, coll := range colls {
		switch cli.Format {
		case "json":
			l, _ := json.Marshal(coll)
			fmt.Println(string(l))
		case "csv":
			fmt.Println(coll.AsCSV())
		}
	}
	return nil
}

type BgpfCLI struct {
	// sub commands
	Projects   ProjectsCmd   `cmd help:"Get information about supported projects"`
	Collectors CollectorsCmd `cmd help:"Get information about supported collectors"`

	// global options
	Format string `help"Output format" default:"json" enum:"json,csv"`

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
	defer os.Stderr.Sync() // flush remaining logs
	handleSignals(ctx, *logp, cancel)

	// calls the appropriate command "Run" method
	err = k.Run(*logp, cliCfg)
	k.FatalIfErrorf(err)
}

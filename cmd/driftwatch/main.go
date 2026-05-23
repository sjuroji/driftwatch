package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/live"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/output"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		manifestPath string
		format       string
		exitOnDrift  bool
	)

	flag.StringVar(&manifestPath, "manifest", "", "path to service manifest YAML file (required)")
	flag.StringVar(&format, "format", "text", "output format: text or json")
	flag.BoolVar(&exitOnDrift, "exit-on-drift", false, "exit with code 2 if drift is detected")
	flag.Parse()

	if manifestPath == "" {
		flag.Usage()
		return fmt.Errorf("--manifest flag is required")
	}

	man, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	fetcher := live.NewStaticFetcher(nil)
	liveStates, err := live.FetchAll(fetcher, []manifest.Manifest{*man})
	if err != nil {
		return fmt.Errorf("fetching live state: %w", err)
	}

	liveState, ok := liveStates[man.Name]
	if !ok {
		return fmt.Errorf("no live state returned for service %q", man.Name)
	}

	report := drift.Detect(*man, liveState)

	formatter, err := output.New(format)
	if err != nil {
		return fmt.Errorf("creating formatter: %w", err)
	}

	if err := formatter.Format(os.Stdout, report); err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	if exitOnDrift && report.HasDrift {
		os.Exit(2)
	}

	return nil
}

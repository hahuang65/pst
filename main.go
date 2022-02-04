package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"git.sr.ht/~hwrd/pst/internal/cli"
	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/tui"
	log "github.com/sirupsen/logrus"
)

func PrintDescription() {
	commandName := path.Base(os.Args[0])

	fmt.Fprintf(flag.CommandLine.Output(),
		`%s is client for https://paste.sr.ht
It can be used in 2 modes: TUI and CLI
CLI mode strictly has functionality for creating pastes with a single file
CLI is the default mode when any arguments and flags are passed in
TUI mode has more functionality to list and view pastes, as well as remove pastes
TUI mode is activated with the -i flag, or when no other flags/arguments are passed in

`, commandName)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", commandName)
}

func main() {
	var (
		debug           bool
		interactiveMode bool
		logLevel        string
		private         bool
		showHelp        bool
		unlisted        bool
	)

	flag.Usage = func() {
		PrintDescription()
		flag.PrintDefaults()
	}
	var cliOpts cli.Opts

	flag.BoolVar(&interactiveMode, "i", false, "interactive mode, automatic if no arguments/flags are passed in")
	flag.BoolVar(&debug, "d", false, "shortcut for -l debug")
	flag.BoolVar(&showHelp, "h", false, "shows this help guide")
	flag.StringVar(&cliOpts.Name, "n", "", "sets a name for the paste")
	flag.StringVar(&logLevel, "l", "warn", "sets loglevel to the specified `level`")
	flag.BoolVar(&private, "p", false, "sets visibility of the paste to private")
	flag.BoolVar(&unlisted, "u", false, "sets visibility of the paste to unlisted")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if debug {
		logLevel = "DEBUG"
	}

	switch strings.ToUpper(logLevel) {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
		log.WithFields(log.Fields{
			"logLevel": logLevel,
		}).Warn("loglevel not recognized. Defaulting to `WARN`")
	}

	if len(os.Args[1:]) == 0 {
		interactiveMode = true
	}

	if interactiveMode {
		log.Debug("Starting in TUI mode")
		tui.Start()
	} else {
		log.Debug("Starting in CLI mode")

		if private {
			cliOpts.Visibility = paste.Private
		} else if unlisted {
			cliOpts.Visibility = paste.Unlisted
		} else {
			cliOpts.Visibility = paste.Public
		}

		cliOpts.Path = flag.Arg(0)

		cli.Start(cliOpts)
	}
}

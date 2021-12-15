package main

import (
	"flag"
	"os"
  "strings"
  
  "git.sr.ht/~hwrd/pst/internal/cli"
  "git.sr.ht/~hwrd/pst/internal/paste"
  "git.sr.ht/~hwrd/pst/internal/tui"
  log "github.com/sirupsen/logrus"
)

func main() {
	var (
		interactiveMode bool
    debug           bool
    logLevel        string
    name            string
		showHelp        bool
    private         bool
    unlisted        bool
    visibility      paste.Visibility
	)

	flag.BoolVar(  &interactiveMode, "i", false,  "run interactively")
	flag.BoolVar(  &debug,           "d", false,  "short for -l debug")
	flag.BoolVar(  &showHelp,        "h", false,  "show help")
	flag.StringVar(&name,            "n", "",     "sets a name for the paste")
	flag.StringVar(&logLevel,        "l", "warn", "sets loglevel to the specified `level`")
	flag.BoolVar(  &private,         "p", false,  "sets visibility of the paste to private")
	flag.BoolVar(  &unlisted,        "u", false,  "sets visibility of the paste to unlisted")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

  if debug {
    logLevel = "DEBUG"
  }

  if private {
    visibility = paste.Private
  } else if unlisted {
    visibility = paste.Unlisted
  } else {
    visibility = paste.Public
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

	if interactiveMode {
    log.Debug("Starting in TUI mode")
    tui.Start()
	} else {
    log.Debug("Starting in CLI mode")
    cli.Start(name, visibility)
	}
}
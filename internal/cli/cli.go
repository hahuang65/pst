package cli

import (
	"bufio"
	"flag"
	"io"
	"os"

	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/util"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
	Name       string
	Path       string
	Visibility paste.Visibility
}

func Start(opts Opts) {
	info, err := os.Stdin.Stat()
	log.WithFields(log.Fields{
		"mode":           info.Mode(),
		"ModeCharDevice": os.ModeCharDevice,
		"size":           info.Size(),
	}).Debug("STDIN info")
	util.CheckError(err)

	var content string

	if (info.Mode() & os.ModeCharDevice) == 0 {
		log.Debug("Receiving input from STDIN")
		var output []rune
		reader := bufio.NewReader(os.Stdin)

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			util.CheckError(err)
			output = append(output, input)
		}

		content = string(output)
	} else if opts.Path != "" {
		log.WithFields(log.Fields{
			"path": opts.Path,
		}).Debug("Reading content from file")

		content, err = util.LoadFile(opts.Path)
		util.CheckError(err)
	} else {
		log.Debug("Did not receive input from STDIN, nor a filepath")

		flag.Usage()
		os.Exit(1)
	}

	paste.Create(opts.Name, opts.Visibility, content)
}

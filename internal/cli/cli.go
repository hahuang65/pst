package cli

import (
	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/util"
)

type Opts struct {
	Name       string
	Path       string
	Visibility paste.Visibility
}

func Start(opts Opts) {
	content, err := util.LoadFile(opts.Path)
	util.CheckError(err)

	paste.Create(opts.Name, opts.Visibility, content)
}

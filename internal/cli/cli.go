package cli

import (
  "flag"
  
  "git.sr.ht/~hwrd/pst/internal/paste"
  "git.sr.ht/~hwrd/pst/internal/util"
)

func Start(name string, visibility paste.Visibility) {
  path := flag.Arg(0)
  content, err := util.LoadFile(path)
  util.CheckError(err)

  paste.Create(name, visibility, content)
}

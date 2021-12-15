package main

import (
  "flag"
)

func startCli(name string, visibility Visibility) {
  path := flag.Arg(0)
  content, err := loadFile(path)
  checkError(err)

  createPaste(name, visibility, content)
}

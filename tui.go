package main

import (
  "flag"
  "fmt"
	"io/ioutil"
  "log"

	tea "github.com/charmbracelet/bubbletea"
)

type errorMsg struct { err error }
type loadedFileMsg string
type model struct {
  content      string
	err          error
  path         string
}

func loadFileCmd(path string) tea.Cmd {
  return func() tea.Msg {
    content, err := loadFile(path)
    if err != nil {
        return errorMsg{err}
    }
    return loadedFileMsg(content)
  }
}

func startTui() {
  // If we're in TUI mode, discard log output
  log.SetOutput(ioutil.Discard)

  var (
    opts    []tea.ProgramOption
    m       model
  )

  m.path = flag.Arg(0)
  p := tea.NewProgram(m, opts...)
  checkError(p.Start())
}

func (m model) Init() tea.Cmd {
	return loadFileCmd(m.path)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case loadedFileMsg:
		m.content = string(msg)
		return m, nil
	case errorMsg:
		m.err = msg.err
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
    return fmt.Sprintf("Error: %s", m.err)
  }
  s := "Interactive Mode:"
	return s + "\n" + m.content
}

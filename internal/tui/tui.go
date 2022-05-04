package tui

import (
	"git.sr.ht/~hwrd/pst/internal/tui/view/list"
	"git.sr.ht/~hwrd/pst/internal/util"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	listView tea.Model
}

func newModel() model {
	return model{
		listView: list.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	newListView, cmd := m.listView.Update(msg)
	m.listView = newListView
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.listView.View()
}

func Start() {
	err := tea.NewProgram(newModel()).Start()
	util.CheckError(err)
}

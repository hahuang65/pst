package tui

import (
	"fmt"

	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/tui/view/list"
	"git.sr.ht/~hwrd/pst/internal/util"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type currentView int

const (
	spinnerView currentView = iota
	listView
)

type model struct {
	listView tea.Model
	spinner  spinner.Model
	view     currentView
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		listView: list.New(),
		spinner:  s,
		view:     spinnerView,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.spinner.Tick,
		m.listView.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.view == spinnerView {
			newSpinner, cmd := m.spinner.Update(msg)
			m.spinner = newSpinner
			cmds = append(cmds, cmd)
		}

	case paste.ListMsg:
		m.view = listView
	}

	newListView, cmd := m.listView.Update(msg)
	m.listView = newListView
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.view == spinnerView {
		return fmt.Sprintf("\n\n   %s Loading pastes\n\n", m.spinner.View())
	} else {
		return m.listView.View()
	}
}

func Start() {
	err := tea.NewProgram(newModel()).Start()
	util.CheckError(err)
}

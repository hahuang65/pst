package tui

import (
	"fmt"

	"git.sr.ht/~hwrd/pst/internal/tui/view"
	"git.sr.ht/~hwrd/pst/internal/tui/view/list"
	"git.sr.ht/~hwrd/pst/internal/util"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	listView    tea.Model
	spinner     spinner.Model
	currentView view.View
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		listView:    list.New(),
		spinner:     s,
		currentView: view.Spinner,
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
		if m.currentView == view.Spinner {
			newSpinner, cmd := m.spinner.Update(msg)
			m.spinner = newSpinner
			cmds = append(cmds, cmd)
		}

	case view.SetMsg:
		m.currentView = view.View(msg)
	}

	newListView, cmd := m.listView.Update(msg)
	m.listView = newListView
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currentView == view.Spinner {
		return fmt.Sprintf("\n\n   %s Loading pastes\n\n", m.spinner.View())
	} else {
		return m.listView.View()
	}
}

func Start() {
	err := tea.NewProgram(newModel()).Start()
	util.CheckError(err)
}

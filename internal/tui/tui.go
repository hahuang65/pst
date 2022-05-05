package tui

import (
	"fmt"

	"git.sr.ht/~hwrd/pst/internal/tui/view"
	"git.sr.ht/~hwrd/pst/internal/tui/view/list"
	"git.sr.ht/~hwrd/pst/internal/tui/view/peek"
	"git.sr.ht/~hwrd/pst/internal/util"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	listView    tea.Model
	peekView    tea.Model
	spinner     spinner.Model
	currentView view.View
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		listView:    list.New(),
		peekView:    peek.New(),
		spinner:     s,
		currentView: view.Spinner,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.listView.Init(),
		m.peekView.Init(),
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

	// Only update the sub-models if it's the currently focused one, or the msg is not a keypress
	// This prevents keypresses in one sub-model from triggering actions in another sub-model
	_, isKeyMsg := msg.(tea.KeyMsg)

	if !isKeyMsg || m.currentView == view.List {
		newListView, cmd := m.listView.Update(msg)
		m.listView = newListView
		cmds = append(cmds, cmd)
	}

	if !isKeyMsg || m.currentView == view.Peek {
		newPeekView, cmd := m.peekView.Update(msg)
		m.peekView = newPeekView
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currentView == view.Spinner {
		return fmt.Sprintf("\n\n   %s Loading pastes\n\n", m.spinner.View())
	} else if m.currentView == view.List {
		return m.listView.View()
	} else if m.currentView == view.Peek {
		return m.peekView.View()
	} else {
		return ""
	}
}

func Start() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	err := p.Start()
	util.CheckError(err)
}

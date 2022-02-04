package tui

import (
	"git.sr.ht/~hwrd/pst/internal/paste"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#036B46", Dark: "#036B46"}).
		Render
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var (
			name string
		)

		if i, ok := m.SelectedItem().(paste.ListItem); ok {
			name = i.Name()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.copy):
				return m.NewStatusMessage(statusMessageStyle("Copied " + name))

			case key.Matches(msg, keys.delete):
				return m.NewStatusMessage(statusMessageStyle("Deleted " + name))

			case key.Matches(msg, keys.open):
				return m.NewStatusMessage(statusMessageStyle("Opened " + name))

			case key.Matches(msg, keys.preview):
				return m.NewStatusMessage(statusMessageStyle("Peeking at " + name))
			}
		}

		return nil
	}

	help := []key.Binding{keys.copy, keys.delete, keys.open, keys.preview}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	copy    key.Binding
	delete  key.Binding
	open    key.Binding
	preview key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.copy,
		d.delete,
		d.open,
		d.preview,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.copy,
			d.delete,
			d.open,
			d.preview,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		copy: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy URL"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete paste"),
		),
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open URL"),
		),
		preview: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "preview paste"),
		),
	}
}

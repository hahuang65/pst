package tui

import (
	"io"

	"git.sr.ht/~hwrd/pst/internal/paste"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
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
			pi paste.ListItem
		)

		if i, ok := m.SelectedItem().(paste.ListItem); ok {
			pi = i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.copy):
				clipboard.WriteAll(pi.URL())
				return m.NewStatusMessage(statusMessageStyle("Copied URL for " + pi.Name()))

			case key.Matches(msg, keys.delete):
				pi.Delete()
				index := m.Index()
				m.RemoveItem(index)
				return m.NewStatusMessage(statusMessageStyle("Deleted " + pi.Name()))

			case key.Matches(msg, keys.open):
				browser.Stdout = io.Discard
				browser.Stderr = io.Discard
				browser.OpenURL(pi.URL())
				return m.NewStatusMessage(statusMessageStyle("Opened " + pi.Name()))

			case key.Matches(msg, keys.preview):
				return m.NewStatusMessage(statusMessageStyle("Peeking at " + pi.Name()))

			case key.Matches(msg, keys.refresh):
				m.SetItems(paste.ListItems(paste.List()))
				return m.NewStatusMessage(statusMessageStyle("Refreshing paste.sr.ht"))
			}
		}

		if len(m.Items()) == 0 {
			keys.delete.SetEnabled(false)
		}
		return nil
	}

	help := []key.Binding{keys.copy, keys.delete, keys.open, keys.preview, keys.refresh}

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
	refresh key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.copy,
		d.delete,
		d.open,
		d.preview,
		d.refresh,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.copy,
			d.delete,
			d.open,
			d.preview,
			d.refresh,
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
			key.WithHelp("d", "delete"),
		),
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open browser"),
		),
		preview: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "preview"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}

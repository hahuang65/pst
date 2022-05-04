package list

import (
	"time"

	"git.sr.ht/~hwrd/pst/internal/paste"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#036B46")).
			Padding(0, 1)
)

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

func New() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	items := paste.ListItems(paste.List())

	delegate := newItemDelegate(delegateKeys)
	pasteList := list.New(items, delegate, 0, 0)
	pasteList.Title = "paste.sr.ht"
	pasteList.Styles.Title = titleStyle
	pasteList.StatusMessageLifetime = time.Second * 5

	return model{
		list:         pasteList,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return paste.List
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := listStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return listStyle.Render(m.list.View())
}

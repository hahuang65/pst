package list

import (
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/tui/view"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hako/durafmt"
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

type item struct {
	paste paste.Paste
}

type fetchPastesMsg []paste.Paste

func New() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	delegate := newItemDelegate(delegateKeys)
	pasteList := list.New([]list.Item{}, delegate, 0, 0)
	pasteList.Title = "paste.sr.ht"
	pasteList.Styles.Title = titleStyle
	pasteList.StatusMessageLifetime = time.Second * 5

	return model{
		list:         pasteList,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return fetchPastes
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := listStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case fetchPastesMsg:
		m.list.SetItems(itemize(msg))
		cmds = append(cmds, view.SetView(view.List))

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

func fetchPastes() tea.Msg {
	return fetchPastesMsg(paste.List())
}

func itemize(ps []paste.Paste) []list.Item {
	list_items := []list.Item{}

	for _, p := range ps {
		list_items = append(list_items, item{paste: p})
	}

	return list_items
}

func (i item) filenames() []string {
	filenames := make([]string, len(i.paste.Files))
	for i, v := range i.paste.Files {
		if v.Filename != "" {
			filenames[i] = v.Filename
		} else {
			filenames[i] = "[UNNAMED]"
		}
	}

	return filenames
}

func (i item) Title() string {
	s := ""

	switch i.paste.Visibility {
	case paste.Public:
		s += "  "
	case paste.Private:
		s += "  "
	case paste.Unlisted:
		s += "  "
	}

	return s + i.Name()
}

func (i item) Name() string {
	return i.paste.Sha
}

func (i item) Description() string {
	return fmt.Sprintf("Files: %s (Created %s ago)",
		strings.Join(i.filenames(), ", "),
		durafmt.ParseShort(time.Since(i.paste.CreatedAt)))
}

func (i item) FilterValue() string {
	// Allow filtering by the SHA or any file names in the paste
	return i.paste.Sha + ", " + strings.Join(i.filenames(), ", ")
}

func (i item) URL() string {
	return i.paste.URL()
}

func (i item) Delete() {
	i.paste.Delete()
}

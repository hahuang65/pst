package peek

import (
	"fmt"
	"strings"

	"git.sr.ht/~hwrd/pst/internal/paste"
	"git.sr.ht/~hwrd/pst/internal/tui/view"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	contents     []paste.Blob
	currentIndex int
	headerTitle  string
	ready        bool
	viewport     viewport.Model
}

type peekBlobMsg struct{}
type peekPasteMsg []paste.Blob

func (m model) Init() tea.Cmd {
	return nil
}

func New() model {
	return model{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" || k == "esc" {
			return m, view.SetView(view.List)
		} else if k == "ctrl+n" {
			m.currentIndex = min(len(m.contents)-1, m.currentIndex+1)
			cmds = append(cmds, peekBlob)
		} else if k == "ctrl+p" {
			m.currentIndex = max(0, m.currentIndex-1)
			cmds = append(cmds, peekBlob)
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

	case peekPasteMsg:
		m.contents = msg
		m.currentIndex = 0
		cmds = append(cmds, peekBlob)

	case peekBlobMsg:
		blob := m.contents[m.currentIndex]
		m.headerTitle = blob.Filename
		m.viewport.SetContent(blob.Contents)
		cmds = append(cmds, view.SetView(view.Peek))
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render(m.headerTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Peek(p paste.Paste) tea.Cmd {
	return func() tea.Msg {
		return peekPasteMsg(p.LoadFiles())
	}
}

func peekBlob() tea.Msg {
	return peekBlobMsg{}
}

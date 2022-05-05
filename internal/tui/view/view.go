package view

import tea "github.com/charmbracelet/bubbletea"

type View int

const (
	Spinner View = iota
	List
)

type SetMsg View

func SetView(v View) tea.Cmd {
	return func() tea.Msg {
		return SetMsg(v)
	}
}

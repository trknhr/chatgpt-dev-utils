package component

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	tea.Model
	Next() (Component, tea.Cmd)
	Prev() (Component, tea.Cmd)
}

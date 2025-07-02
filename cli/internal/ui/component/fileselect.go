package component

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type FileSelect struct {
	Title     string
	Viewport  viewport.Model
	Message   string
	Width     int
	Height    int
	Cursor    int
	FlatFiles []*file.FileNode
	Selected  []*file.FileNode
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 0).
			Margin(0, 0)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))
)

func NewFileSelect(flat []*file.FileNode, selected []*file.FileNode, vp viewport.Model, cursor, w, h int, msg string) *FileSelect {
	return &FileSelect{
		Title:     "Step 2/4: Select Files",
		FlatFiles: flat,
		Selected:  selected,
		Viewport:  vp,
		Cursor:    cursor,
		Width:     w,
		Height:    h,
		Message:   msg,
	}
}

func (f *FileSelect) Render() string {
	title := titleStyle.Render(f.Title)
	help := helpStyle.Render("[↑↓ Navigate] [Enter: Toggle folder] [Space: Select file] [Tab: Next]")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		f.Viewport.View(),
		"",
		help,
		f.Message,
	)

	return content
}

func (f *FileSelect) Init() tea.Cmd { return nil }

func (m *FileSelect) Next() (Component, tea.Cmd) {
	return nil, nil
}

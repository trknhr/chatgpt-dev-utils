package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type PromptTypeModel struct {
	cursor        int
	width, height int
}

func NewPromptType(w, h int) *PromptTypeModel { return &PromptTypeModel{width: w, height: h} }

func (m PromptTypeModel) Init() tea.Cmd { return nil }

func (m PromptTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = 0
		case "down", "j":
			m.cursor = 1
		}
	}
	return m, nil
}

func (m PromptTypeModel) View() string {
	options := []string{
		"File based Prompt",
		"Git based Prompt",
	}

	content := ""
	for i, option := range options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			option = selectedStyle.Render(option)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, option)
	}

	return RenderLayout(
		"Step 1: Choose Prompt Type",
		content,
		"[↑↓ Navigate] [Tab: Next] [Ctrl+C: Quit]",
		m.width,
		m.height,
	)
}

func (m PromptTypeModel) Next() (Component, tea.Cmd) {
	if m.cursor == 0 {
		// File selection path
		root := file.BuildFileTree(".")
		flat := file.FlattenFileTree(root)

		// Create viewport with proper dimensions
		headerHeight := 4
		footerHeight := 4
		viewportHeight := m.height - headerHeight - footerHeight
		if viewportHeight < 3 {
			viewportHeight = 3
		}
		viewportWidth := m.width - 4
		if viewportWidth < 20 {
			viewportWidth = 20
		}

		vp := viewport.New(viewportWidth, viewportHeight)
		// Initialize content
		content := ""
		for i, node := range flat {
			prefix := " "
			if i == 0 {
				prefix = ">"
			}
			content += fmt.Sprintf("%s %s\n", prefix, file.RenderFileNode(node))
		}
		vp.SetContent(content)

		selectPage := NewFileSelect(flat, []*file.FileNode{}, vp, 0, m.width, m.height, "")
		return selectPage, nil
	}

	// Git template selection path
	templates := []string{"Code Review", "Commit Message", "Change Summary", "Custom..."}
	return NewTemplateSelect("git", templates, nil, m.width, m.height), nil
}
func (m PromptTypeModel) Prev() (Component, tea.Cmd) { return m, nil }

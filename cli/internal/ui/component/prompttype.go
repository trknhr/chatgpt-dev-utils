package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
	"github.com/trknhr/chatgpt-dev-utils/internal/ui/component"
)

type PromptTypeModel struct {
	cursor        int
	width, height int
}

func NewPromptType(w, h int) *PromptTypeModel { return &PromptTypeModel{width: w, height: h} }

func (m PromptTypeModel) Init() tea.Cmd { return nil }

func (m PromptTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		switch k.String() {
		case "up", "k":
			m.cursor = 0
		case "down", "j":
			m.cursor = 1
		}
	}
	return m, nil
}

func (m PromptTypeModel) View() string { /* your previous viewPromptType */ }

func (m PromptTypeModel) Next() (Component, tea.Cmd) {
	if m.cursor == 0 {
		root := file.BuildFileTree(".")
		flat := file.FlattenFileTree(root)

		vp := viewport.New(m.width-4, m.height-8)
		// 最初の内容を詰めておく
		content := ""
		for i, node := range flat {
			prefix := " "
			if i == 0 {
				prefix = ">"
			}
			content += fmt.Sprintf("%s %s\n", prefix, file.RenderFileNode(node))
		}
		vp.SetContent(content)

		selectPage := component.NewFileSelect(flat, []*file.FileNode{}, vp, 0, m.width, m.height, "")
		return selectPage, nil
	}

	return NewGitTemplate(m.width, m.height), nil
}
func (m PromptTypeModel) Prev() (Component, tea.Cmd) { return m, nil }

package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestPromptTypeModel(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewPromptType creates model with correct dimensions",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				assert.Equal(t, 80, model.width)
				assert.Equal(t, 24, model.height)
				assert.Equal(t, 0, model.cursor)
			},
		},
		{
			name: "Init returns nil command",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				cmd := model.Init()
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Update handles up key",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				model.cursor = 1

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
				newModel, cmd := model.Update(msg)
				m := newModel.(PromptTypeModel)

				assert.Equal(t, 0, m.cursor)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Update handles down key",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				model.cursor = 0

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
				newModel, cmd := model.Update(msg)
				m := newModel.(PromptTypeModel)

				assert.Equal(t, 1, m.cursor)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, cmd := model.Update(msg)
				m := newModel.(PromptTypeModel)

				assert.Equal(t, 100, m.width)
				assert.Equal(t, 30, m.height)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "View renders correctly",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				view := model.View()

				assert.Contains(t, view, "Step 1: Choose Prompt Type")
				assert.Contains(t, view, "File based Prompt")
				assert.Contains(t, view, "Git based Prompt")
				assert.Contains(t, view, "[↑↓ Navigate] [Tab: Next] [Ctrl+C: Quit]")
			},
		},
		{
			name: "Next returns FileSelect when cursor is 0",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				model.cursor = 0

				next, cmd := model.Next()
				_, ok := next.(*FileSelect)

				assert.True(t, ok)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Next returns TemplateSelect for git when cursor is 1",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)
				model.cursor = 1

				next, cmd := model.Next()
				ts, ok := next.(*TemplateSelect)

				assert.True(t, ok)
				assert.Equal(t, "git", ts.PromptType)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Prev returns self",
			test: func(t *testing.T) {
				model := NewPromptType(80, 24)

				prev, cmd := model.Prev()
				m, ok := prev.(PromptTypeModel)

				assert.True(t, ok)
				assert.Equal(t, *model, m)
				assert.Nil(t, cmd)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

func TestEdit(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewEdit creates model correctly",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				edit := NewEdit("git", "Code Review", "Test template", files, 80, 24)

				assert.Equal(t, "git", edit.PromptType)
				assert.Equal(t, "Code Review", edit.SelectedTemplate)
				assert.Equal(t, "Test template", edit.TemplateContent)
				assert.Equal(t, files, edit.SelectedFiles)
				assert.Equal(t, 80, edit.Width)
				assert.Equal(t, 24, edit.Height)
				assert.Equal(t, "Test template", edit.Textarea.Value())
			},
		},
		{
			name: "Init returns blink command",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Template", "Content", nil, 80, 24)
				cmd := edit.Init()
				assert.NotNil(t, cmd)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Template", "Content", nil, 80, 24)

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, _ := edit.Update(msg)
				updated := newModel.(*Edit)

				assert.Equal(t, 100, updated.Width)
				assert.Equal(t, 30, updated.Height)
			},
		},
		{
			name: "Update handles tab key to move to next",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Template", "Content", nil, 80, 24)

				msg := tea.KeyMsg{Type: tea.KeyTab}
				newModel, _ := edit.Update(msg)

				// Should still be Edit component since Tab is handled specially
				_, ok := newModel.(*Edit)
				assert.True(t, ok)
			},
		},
		{
			name: "View renders correctly for git",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Code Review", "Content", nil, 80, 24)
				view := edit.View()

				assert.Contains(t, view, "Step 3: Review & Edit")
				assert.Contains(t, view, "[↑↓←→ Type freely] [Tab: Next] [Esc: Back]")
			},
		},
		{
			name: "View renders correctly for file",
			test: func(t *testing.T) {
				edit := NewEdit("file", "Documentation", "Content", nil, 80, 24)
				view := edit.View()

				assert.Contains(t, view, "Step 4: Review & Edit")
			},
		},
		{
			name: "View handles small terminal",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Template", "Content", nil, 15, 8)
				view := edit.View()

				assert.Equal(t, "Your terminal is too small.", view)
			},
		},
		{
			name: "Next returns Final component",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				edit := NewEdit("git", "Code Review", "Modified content", files, 80, 24)
				edit.Textarea.SetValue("Final prompt content")

				next, _ := edit.Next()
				final, ok := next.(*Final)

				assert.True(t, ok)
				assert.Equal(t, "git", final.PromptType)
				assert.Equal(t, "Code Review", final.SelectedTemplate)
				assert.Equal(t, "Final prompt content", final.FinalPrompt)
				assert.Equal(t, files, final.SelectedFiles)
			},
		},
		{
			name: "Prev returns TemplateSelect for git",
			test: func(t *testing.T) {
				edit := NewEdit("git", "Code Review", "Content", nil, 80, 24)

				prev, _ := edit.Prev()
				ts, ok := prev.(*TemplateSelect)

				assert.True(t, ok)
				assert.Equal(t, "git", ts.PromptType)
				assert.Contains(t, ts.Templates, "Code Review")
			},
		},
		{
			name: "Prev returns TemplateSelect for file",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				edit := NewEdit("file", "Documentation", "Content", files, 80, 24)

				prev, _ := edit.Prev()
				ts, ok := prev.(*TemplateSelect)

				assert.True(t, ok)
				assert.Equal(t, "file", ts.PromptType)
				assert.Contains(t, ts.Templates, "Documentation")
				assert.Equal(t, files, ts.SelectedFiles)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
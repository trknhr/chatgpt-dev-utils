package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

func TestTemplateSelect(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewTemplateSelect creates model correctly",
			test: func(t *testing.T) {
				templates := []string{"Template1", "Template2"}
				files := []*file.FileNode{{Path: "test.go"}}
				ts := NewTemplateSelect("file", templates, files, 80, 24)

				assert.Equal(t, "file", ts.PromptType)
				assert.Equal(t, templates, ts.Templates)
				assert.Equal(t, files, ts.SelectedFiles)
				assert.Equal(t, 0, ts.Cursor)
				assert.Equal(t, 80, ts.Width)
				assert.Equal(t, 24, ts.Height)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"Template1"}, nil, 80, 24)

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, _ := ts.Update(msg)
				updated := newModel.(*TemplateSelect)

				assert.Equal(t, 100, updated.Width)
				assert.Equal(t, 30, updated.Height)
			},
		},
		{
			name: "Update handles up navigation",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"T1", "T2"}, nil, 80, 24)
				ts.Cursor = 1

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
				newModel, _ := ts.Update(msg)
				updated := newModel.(*TemplateSelect)

				assert.Equal(t, 0, updated.Cursor)
			},
		},
		{
			name: "Update handles down navigation",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"T1", "T2"}, nil, 80, 24)

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
				newModel, _ := ts.Update(msg)
				updated := newModel.(*TemplateSelect)

				assert.Equal(t, 1, updated.Cursor)
			},
		},
		{
			name: "View renders git template title correctly",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"Code Review"}, nil, 80, 24)
				view := ts.View()

				assert.Contains(t, view, "Step 2: Choose Prompt Template")
				assert.Contains(t, view, "Code Review")
			},
		},
		{
			name: "View renders file template title correctly",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("file", []string{"Documentation"}, nil, 80, 24)
				view := ts.View()

				assert.Contains(t, view, "Step 3: Choose Prompt Template")
				assert.Contains(t, view, "Documentation")
			},
		},
		{
			name: "Next returns Edit component for git",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"Code Review", "Commit Message"}, nil, 80, 24)
				ts.Cursor = 1 // Select "Commit Message"

				next, _ := ts.Next()
				edit, ok := next.(*Edit)

				assert.True(t, ok)
				assert.Equal(t, "git", edit.PromptType)
				assert.Equal(t, "Commit Message", edit.SelectedTemplate)
			},
		},
		{
			name: "Next returns Edit component for file",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				ts := NewTemplateSelect("file", []string{"Code Review", "Documentation"}, files, 80, 24)
				ts.Cursor = 0 // Select "Code Review"

				next, _ := ts.Next()
				edit, ok := next.(*Edit)

				assert.True(t, ok)
				assert.Equal(t, "file", edit.PromptType)
				assert.Equal(t, "Code Review", edit.SelectedTemplate)
				assert.Equal(t, files, edit.SelectedFiles)
			},
		},
		{
			name: "Prev returns FileSelect for file type",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				ts := NewTemplateSelect("file", []string{"Template1"}, files, 80, 24)

				prev, _ := ts.Prev()
				_, ok := prev.(*FileSelect)

				assert.True(t, ok)
			},
		},
		{
			name: "Prev returns PromptType for git type",
			test: func(t *testing.T) {
				ts := NewTemplateSelect("git", []string{"Template1"}, nil, 80, 24)

				prev, _ := ts.Prev()
				_, ok := prev.(*PromptTypeModel)

				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
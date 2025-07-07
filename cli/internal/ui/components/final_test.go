package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

func TestFinal(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewFinal creates model correctly",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				broadcastChan := make(chan<- string)
				clientsCount := func() int { return 1 }

				final := NewFinal("git", "Code Review", "Final prompt", files, 80, 24, true, broadcastChan, clientsCount)

				assert.Equal(t, "git", final.PromptType)
				assert.Equal(t, "Code Review", final.SelectedTemplate)
				assert.Equal(t, "Final prompt", final.FinalPrompt)
				assert.Equal(t, files, final.SelectedFiles)
				assert.Equal(t, 80, final.Width)
				assert.Equal(t, 24, final.Height)
				assert.True(t, final.ExtensionConnected)
				assert.NotNil(t, final.BroadcastChan)
				assert.NotNil(t, final.ClientsCount)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				final := NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, nil)

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, _ := final.Update(msg)
				updated := newModel.(*Final)

				assert.Equal(t, 100, updated.Width)
				assert.Equal(t, 30, updated.Height)
			},
		},
		{
			name: "Update handles CheckConnectionMsg",
			test: func(t *testing.T) {
				clientsCount := func() int { return 2 }
				final := NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, clientsCount)

				msg := CheckConnectionMsg{}
				newModel, _ := final.Update(msg)
				updated := newModel.(*Final)

				assert.True(t, updated.ExtensionConnected)
			},
		},
		{
			name: "Update handles copy command",
			test: func(t *testing.T) {
				final := NewFinal("git", "Template", "Test prompt", nil, 80, 24, false, nil, nil)

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")}
				newModel, _ := final.Update(msg)
				updated := newModel.(*Final)

				assert.Equal(t, "Copied to clipboard!", updated.Message)
			},
		},
		{
			name: "View renders git prompt correctly",
			test: func(t *testing.T) {
				final := NewFinal("git", "Code Review", "Review this code", nil, 80, 24, false, nil, nil)
				view := final.View()

				assert.Contains(t, view, "Step 4: Copy Prompt")
				assert.Contains(t, view, "Ready to copy:")
				assert.Contains(t, view, "[C: Copy with Content] [Esc: Back]")
			},
		},
		{
			name: "View renders file prompt correctly",
			test: func(t *testing.T) {
				files := []*file.FileNode{
					{Path: "test1.go"},
					{Path: "test2.go"},
				}
				final := NewFinal("file", "Documentation", "Document these files", files, 80, 24, false, nil, nil)
				view := final.View()

				assert.Contains(t, view, "Step 5: Copy Prompt")
				assert.Contains(t, view, "Selected files:")
				assert.Contains(t, view, "test1.go")
				assert.Contains(t, view, "test2.go")
			},
		},
		{
			name: "View shows extension option when connected",
			test: func(t *testing.T) {
				final := NewFinal("git", "Template", "Prompt", nil, 80, 24, true, nil, nil)
				view := final.View()

				assert.Contains(t, view, "[E: Send to Extension]")
			},
		},
		{
			name: "View shows message when present",
			test: func(t *testing.T) {
				final := NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, nil)
				final.Message = "Test message"
				view := final.View()

				assert.Contains(t, view, "Test message")
			},
		},
		{
			name: "Next returns self",
			test: func(t *testing.T) {
				final := NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, nil)

				next, _ := final.Next()
				f, ok := next.(*Final)

				assert.True(t, ok)
				assert.Equal(t, final, f)
			},
		},
		{
			name: "Prev returns Edit for git",
			test: func(t *testing.T) {
				final := NewFinal("git", "Code Review", "Prompt", nil, 80, 24, false, nil, nil)

				prev, _ := final.Prev()
				edit, ok := prev.(*Edit)

				assert.True(t, ok)
				assert.Equal(t, "git", edit.PromptType)
				assert.Equal(t, "Code Review", edit.SelectedTemplate)
			},
		},
		{
			name: "Prev returns Edit for file",
			test: func(t *testing.T) {
				files := []*file.FileNode{{Path: "test.go"}}
				final := NewFinal("file", "Documentation", "Prompt", files, 80, 24, false, nil, nil)

				prev, _ := final.Prev()
				edit, ok := prev.(*Edit)

				assert.True(t, ok)
				assert.Equal(t, "file", edit.PromptType)
				assert.Equal(t, "Documentation", edit.SelectedTemplate)
				assert.Equal(t, files, edit.SelectedFiles)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
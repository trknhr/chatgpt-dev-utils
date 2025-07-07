package components

import (
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

func TestFileSelect(t *testing.T) {
	// Create test file nodes
	createTestFileNodes := func() []*file.FileNode {
		root := &file.FileNode{
			Name:   "test",
			Path:   "test",
			IsDir:  true,
			IsOpen: true,
		}
		file1 := &file.FileNode{
			Name:   "file1.go",
			Path:   "test/file1.go",
			IsDir:  false,
			Parent: root,
		}
		file2 := &file.FileNode{
			Name:   "file2.go",
			Path:   "test/file2.go",
			IsDir:  false,
			Parent: root,
		}
		root.Children = []*file.FileNode{file1, file2}
		return []*file.FileNode{root, file1, file2}
	}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewFileSelect creates model correctly",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "test message")

				assert.Equal(t, "Step 2: Select Files", fs.Title)
				assert.Equal(t, flat, fs.FlatFiles)
				assert.Equal(t, 0, fs.Cursor)
				assert.Equal(t, 80, fs.Width)
				assert.Equal(t, 24, fs.Height)
				assert.Equal(t, "test message", fs.Message)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "")

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, _ := fs.Update(msg)
				updated := newModel.(*FileSelect)

				assert.Equal(t, 100, updated.Width)
				assert.Equal(t, 30, updated.Height)
				assert.Equal(t, 96, updated.Viewport.Width) // 100 - 4
				assert.Equal(t, 22, updated.Viewport.Height) // 30 - 8
			},
		},
		{
			name: "Update handles up navigation",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 1, 80, 24, "")

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
				newModel, _ := fs.Update(msg)
				updated := newModel.(*FileSelect)

				assert.Equal(t, 0, updated.Cursor)
			},
		},
		{
			name: "Update handles down navigation",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "")

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
				newModel, _ := fs.Update(msg)
				updated := newModel.(*FileSelect)

				assert.Equal(t, 1, updated.Cursor)
			},
		},
		{
			name: "Update handles file selection with space",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 1, 80, 24, "") // cursor on file1

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")}
				newModel, _ := fs.Update(msg)
				updated := newModel.(*FileSelect)

				assert.Len(t, updated.Selected, 1)
				assert.True(t, updated.FlatFiles[1].Selected)
			},
		},
		{
			name: "Update toggles folder open/close with enter",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "") // cursor on root dir

				msg := tea.KeyMsg{Type: tea.KeyEnter}
				_, _ = fs.Update(msg)

				// After toggling, the folder should be closed
				// The original node's state is modified
				assert.False(t, flat[0].IsOpen)
			},
		},
		{
			name: "View renders correctly",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "")

				view := fs.View()
				assert.Contains(t, view, "Step 2: Select Files")
				assert.Contains(t, view, "[↑↓ Navigate] [Enter: Toggle folder] [Space: Select file] [Tab: Next]")
			},
		},
		{
			name: "Next returns TemplateSelect when files selected",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, []*file.FileNode{flat[1]}, vp, 0, 80, 24, "")

				next, _ := fs.Next()
				ts, ok := next.(*TemplateSelect)

				assert.True(t, ok)
				assert.Equal(t, "file", ts.PromptType)
				assert.Len(t, ts.SelectedFiles, 1)
			},
		},
		{
			name: "Next returns self when no files selected",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "")

				next, _ := fs.Next()
				_, ok := next.(*FileSelect)

				assert.True(t, ok)
			},
		},
		{
			name: "Prev returns PromptType",
			test: func(t *testing.T) {
				flat := createTestFileNodes()
				vp := viewport.New(80, 20)
				fs := NewFileSelect(flat, nil, vp, 0, 80, 24, "")

				prev, _ := fs.Prev()
				_, ok := prev.(*PromptTypeModel)

				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
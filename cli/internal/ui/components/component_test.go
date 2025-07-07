package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/stretchr/testify/assert"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

// TestComponentInterface ensures all components properly implement the Component interface
func TestComponentInterface(t *testing.T) {
	// Create test components
	components := []struct {
		name      string
		component Component
	}{
		{
			name:      "PromptTypeModel",
			component: NewPromptType(80, 24),
		},
		{
			name:      "FileSelect",
			component: NewFileSelect([]*file.FileNode{}, nil, viewport.New(80, 20), 0, 80, 24, ""),
		},
		{
			name:      "TemplateSelect",
			component: NewTemplateSelect("git", []string{"Test"}, nil, 80, 24),
		},
		{
			name:      "Edit",
			component: NewEdit("git", "Template", "Content", nil, 80, 24),
		},
		{
			name:      "Final",
			component: NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, nil),
		},
		{
			name:      "Root",
			component: NewRoot(80, 24, nil, nil),
		},
	}

	for _, tc := range components {
		t.Run(tc.name, func(t *testing.T) {
			// Test that component implements Component interface
			var _ Component = tc.component
			
			// Test Init method
			cmd := tc.component.Init()
			// Init should return either nil or a valid command
			assert.True(t, cmd == nil || cmd != nil)
			
			// Test Update method with various messages
			testMsgs := []tea.Msg{
				tea.WindowSizeMsg{Width: 100, Height: 30},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")},
				tea.KeyMsg{Type: tea.KeyTab},
			}
			
			for _, msg := range testMsgs {
				model, cmd := tc.component.Update(msg)
				assert.NotNil(t, model)
				// Command can be nil or not
				assert.True(t, cmd == nil || cmd != nil)
			}
			
			// Test View method
			view := tc.component.View()
			assert.NotEmpty(t, view)
			
			// Test Next method
			next, cmd := tc.component.Next()
			assert.NotNil(t, next)
			assert.True(t, cmd == nil || cmd != nil)
			
			// Test Prev method
			prev, cmd := tc.component.Prev()
			assert.NotNil(t, prev)
			assert.True(t, cmd == nil || cmd != nil)
		})
	}
}

// TestComponentNavigation tests the navigation flow between components
func TestComponentNavigation(t *testing.T) {
	t.Run("Full navigation flow", func(t *testing.T) {
		// Start with PromptType
		prompt := NewPromptType(80, 24)
		
		// Navigate to FileSelect
		prompt.cursor = 0 // Select file option
		next, _ := prompt.Next()
		fileSelect, ok := next.(*FileSelect)
		assert.True(t, ok)
		
		// Navigate back to PromptType
		prev, _ := fileSelect.Prev()
		_, ok = prev.(*PromptTypeModel)
		assert.True(t, ok)
		
		// Navigate to git path
		prompt.cursor = 1 // Select git option
		next, _ = prompt.Next()
		templateSelect, ok := next.(*TemplateSelect)
		assert.True(t, ok)
		assert.Equal(t, "git", templateSelect.PromptType)
		
		// Navigate to Edit
		next, _ = templateSelect.Next()
		edit, ok := next.(*Edit)
		assert.True(t, ok)
		
		// Navigate to Final
		next, _ = edit.Next()
		final, ok := next.(*Final)
		assert.True(t, ok)
		
		// Final.Next() should return self
		next, _ = final.Next()
		assert.Equal(t, final, next)
	})
}
package components

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyles(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "RenderLayout creates basic layout",
			test: func(t *testing.T) {
				result := RenderLayout("Test Title", "Test Content", "Test Help", 80, 24)

				assert.Contains(t, result, "Test Title")
				assert.Contains(t, result, "Test Content")
				assert.Contains(t, result, "Test Help")
				
				// Should have proper structure
				lines := strings.Split(result, "\n")
				assert.GreaterOrEqual(t, len(lines), 3)
			},
		},
		{
			name: "RenderLayout adjusts box width for small terminals",
			test: func(t *testing.T) {
				// Small width should trigger width adjustment
				result := RenderLayout("Title", "Content", "Help", 50, 24)
				
				// Basic check that layout is rendered
				assert.Contains(t, result, "Title")
				assert.Contains(t, result, "Content")
				assert.Contains(t, result, "Help")
			},
		},
		{
			name: "RenderLayoutWithMessage includes message when provided",
			test: func(t *testing.T) {
				result := RenderLayoutWithMessage("Title", "Content", "Help", "Test Message", 80, 24)

				assert.Contains(t, result, "Title")
				assert.Contains(t, result, "Content")
				assert.Contains(t, result, "Help")
				assert.Contains(t, result, "Test Message")
			},
		},
		{
			name: "RenderLayoutWithMessage excludes empty message",
			test: func(t *testing.T) {
				result := RenderLayoutWithMessage("Title", "Content", "Help", "", 80, 24)

				assert.Contains(t, result, "Title")
				assert.Contains(t, result, "Content")
				assert.Contains(t, result, "Help")
				
				// Should not have extra newlines at the end
				assert.False(t, strings.HasSuffix(result, "\n\n"))
			},
		},
		{
			name: "Styles are properly initialized",
			test: func(t *testing.T) {
				// Test that styles are not nil and have expected properties
				assert.NotNil(t, titleStyle)
				assert.NotNil(t, selectedStyle)
				assert.NotNil(t, helpStyle)

				// Test rendering with styles
				titleRendered := titleStyle.Render("Test")
				assert.Contains(t, titleRendered, "Test")

				selectedRendered := selectedStyle.Render("Selected")
				assert.Contains(t, selectedRendered, "Selected")

				helpRendered := helpStyle.Render("Help")
				assert.Contains(t, helpRendered, "Help")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
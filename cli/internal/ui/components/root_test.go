package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewRoot creates root with PromptType child",
			test: func(t *testing.T) {
				broadcastChan := make(chan<- string)
				clientsCount := func() int { return 0 }
				root := NewRoot(80, 24, broadcastChan, clientsCount)

				assert.Equal(t, 80, root.width)
				assert.Equal(t, 24, root.height)
				assert.NotNil(t, root.child)
				assert.NotNil(t, root.broadcastChan)
				assert.NotNil(t, root.clientsCount)

				_, ok := root.child.(*PromptTypeModel)
				assert.True(t, ok)
			},
		},
		{
			name: "Init delegates to child",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				cmd := root.Init()
				// Child's Init returns nil, so root should too
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Update handles window resize",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)

				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				newModel, _ := root.Update(msg)
				updated := newModel.(*Root)

				assert.Equal(t, 100, updated.width)
				assert.Equal(t, 30, updated.height)
			},
		},
		{
			name: "Update handles CheckConnectionMsg",
			test: func(t *testing.T) {
				clientsCount := func() int { return 3 }
				root := NewRoot(80, 24, nil, clientsCount)

				msg := CheckConnectionMsg{}
				newModel, _ := root.Update(msg)
				updated := newModel.(*Root)

				assert.True(t, updated.extensionConnected)
			},
		},
		{
			name: "Update handles tab navigation forward",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				// Set up root with a mock component that has Next
				mockChild := &mockComponent{
					nextComponent: NewPromptType(80, 24),
				}
				root.child = mockChild

				msg := tea.KeyMsg{Type: tea.KeyTab}
				newModel, _ := root.Update(msg)
				updated := newModel.(*Root)

				// Should have navigated to the next component
				_, ok := updated.child.(*PromptTypeModel)
				assert.True(t, ok)
			},
		},
		{
			name: "Update handles esc navigation backward",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				// Set up root with a mock component that has Prev
				mockChild := &mockComponent{
					prevComponent: NewPromptType(80, 24),
				}
				root.child = mockChild

				msg := tea.KeyMsg{Type: tea.KeyEsc}
				newModel, _ := root.Update(msg)
				updated := newModel.(*Root)

				// Should have navigated to the prev component
				_, ok := updated.child.(*PromptTypeModel)
				assert.True(t, ok)
			},
		},
		{
			name: "Update injects WebSocket context to Final component",
			test: func(t *testing.T) {
				broadcastChan := make(chan<- string)
				clientsCount := func() int { return 1 }
				root := NewRoot(80, 24, broadcastChan, clientsCount)
				root.extensionConnected = true

				// Set up with mock that returns Final component
				finalComp := NewFinal("git", "Template", "Prompt", nil, 80, 24, false, nil, nil)
				mockChild := &mockComponent{
					nextComponent: finalComp,
				}
				root.child = mockChild

				msg := tea.KeyMsg{Type: tea.KeyTab}
				newModel, _ := root.Update(msg)
				updated := newModel.(*Root)

				// Check if Final component got the WebSocket context
				final, ok := updated.child.(*Final)
				assert.True(t, ok)
				assert.True(t, final.ExtensionConnected)
				assert.NotNil(t, final.BroadcastChan)
				assert.NotNil(t, final.ClientsCount)
			},
		},
		{
			name: "View delegates to child",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				view := root.View()

				// Should contain content from PromptType view
				assert.Contains(t, view, "Step 1: Choose Prompt Type")
			},
		},
		{
			name: "Next returns self",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				next, cmd := root.Next()

				assert.Equal(t, root, next)
				assert.Nil(t, cmd)
			},
		},
		{
			name: "Prev returns self",
			test: func(t *testing.T) {
				root := NewRoot(80, 24, nil, nil)
				prev, cmd := root.Prev()

				assert.Equal(t, root, prev)
				assert.Nil(t, cmd)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// mockComponent helps test navigation
type mockComponent struct {
	nextComponent Component
	prevComponent Component
}

func (m *mockComponent) Init() tea.Cmd { return nil }
func (m *mockComponent) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *mockComponent) View() string { return "mock view" }
func (m *mockComponent) Next() (Component, tea.Cmd) { return m.nextComponent, nil }
func (m *mockComponent) Prev() (Component, tea.Cmd) { return m.prevComponent, nil }
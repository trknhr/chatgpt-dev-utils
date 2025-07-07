package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trknhr/chatgpt-dev-utils/internal/ui/components"
)

// Model is the main UI model that orchestrates components
type Model struct {
	root               *components.Root
	broadcastChan      chan<- string
	clientsCount       func() int
}

// InitialModel creates the initial TUI model using components
func InitialModel(broadcastChan chan<- string, clientsCount func() int) Model {
	// Get initial terminal size
	width, height := 80, 24 // default size
	
	// Create root component with WebSocket context
	root := components.NewRoot(width, height, broadcastChan, clientsCount)
	
	return Model{
		root:          root,
		broadcastChan: broadcastChan,
		clientsCount:  clientsCount,
	}
}

func (m Model) Init() tea.Cmd {
	// Start connection check timer
	return tea.Batch(
		m.root.Init(),
		tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return components.CheckConnectionMsg{}
		}),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case components.CheckConnectionMsg:
		// Pass to root and restart timer
		updatedRoot, cmd := m.root.Update(msg)
		m.root = updatedRoot.(*components.Root)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		// Restart the timer
		cmds = append(cmds, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return components.CheckConnectionMsg{}
		}))
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Handle global keys
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// Delegate all other updates to root component
	updatedRoot, cmd := m.root.Update(msg)
	m.root = updatedRoot.(*components.Root)
	
	return m, cmd
}

func (m Model) View() string {
	return m.root.View()
}
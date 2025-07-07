package components

import tea "github.com/charmbracelet/bubbletea"

type Root struct {
	child              Component
	width, height      int
	broadcastChan      chan<- string
	clientsCount       func() int
	extensionConnected bool
}

func NewRoot(w, h int, broadcastChan chan<- string, clientsCount func() int) *Root {
	return &Root{
		child:         NewPromptType(w, h),
		width:         w,
		height:        h,
		broadcastChan: broadcastChan,
		clientsCount:  clientsCount,
	}
}

func (r *Root) Init() tea.Cmd { return r.child.Init() }

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width, r.height = msg.Width, msg.Height
		// Let the message propagate to child components
		
	case CheckConnectionMsg:
		// Update extension connection status
		if r.clientsCount != nil {
			r.extensionConnected = r.clientsCount() > 0
		}
		return r, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Navigate forward
			nextChild, cmd := r.child.Next()
			if nextChild != nil {
				// Pass WebSocket context to Final component
				if final, ok := nextChild.(*Final); ok {
					final.ExtensionConnected = r.extensionConnected
					final.BroadcastChan = r.broadcastChan
					final.ClientsCount = r.clientsCount
				}
				r.child = nextChild
			}
			return r, cmd
		case "esc":
			// Navigate backward
			prevChild, cmd := r.child.Prev()
			if prevChild != nil {
				r.child = prevChild
			}
			return r, cmd
		}
	}

	// Delegate to child component
	updated, cmd := r.child.Update(msg)
	if updated != nil {
		r.child = updated.(Component)
	}
	return r, cmd
}

func (r *Root) View() string { return r.child.View() }

func (r *Root) Next() (Component, tea.Cmd) {
	// Root doesn't navigate, it manages child navigation
	return r, nil
}

func (r *Root) Prev() (Component, tea.Cmd) {
	// Root doesn't navigate, it manages child navigation
	return r, nil
}

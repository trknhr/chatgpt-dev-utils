package component

import tea "github.com/charmbracelet/bubbletea"

type Root struct {
	child         Component
	width, height int
}

func NewRoot(w, h int) *Root { return &Root{child: NewPromptType(w, h), width: w, height: h} }

func (r *Root) Init() tea.Cmd { return r.child.Init() }

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		r.width, r.height = msg.(tea.WindowSizeMsg).Width, msg.(tea.WindowSizeMsg).Height
	}

	updated, cmd := r.child.Update(msg)
	r.child = updated.(Component)

	switch msg.(type) {
	case nextMsg:
		r.child, _ = r.child.Next()
	case prevMsg:
		r.child, _ = r.child.Prev()
	}
	return r, cmd
}

func (r *Root) View() string { return r.child.View() }

package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type FileSelect struct {
	Title     string
	Viewport  viewport.Model
	Message   string
	Width     int
	Height    int
	Cursor    int
	FlatFiles []*file.FileNode
	Selected  []*file.FileNode
}

func NewFileSelect(flat []*file.FileNode, selected []*file.FileNode, vp viewport.Model, cursor, w, h int, msg string) *FileSelect {
	return &FileSelect{
		Title:     "Step 2: Select Files",
		FlatFiles: flat,
		Selected:  selected,
		Viewport:  vp,
		Cursor:    cursor,
		Width:     w,
		Height:    h,
		Message:   msg,
	}
}

func (f *FileSelect) Init() tea.Cmd { return nil }

func (f *FileSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.Width = msg.Width
		f.Height = msg.Height
		// Update viewport dimensions
		headerHeight := 4
		footerHeight := 4
		viewportHeight := f.Height - headerHeight - footerHeight
		if viewportHeight < 3 {
			viewportHeight = 3
		}
		viewportWidth := f.Width - 4
		if viewportWidth < 20 {
			viewportWidth = 20
		}
		f.Viewport.Width = viewportWidth
		f.Viewport.Height = viewportHeight

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if f.Cursor > 0 {
				f.Cursor--
				f.updateViewportContent()
				f.ensureCursorVisible()
			}
		case "down", "j":
			if f.Cursor < len(f.FlatFiles)-1 {
				f.Cursor++
				f.updateViewportContent()
				f.ensureCursorVisible()
			}
		case "enter":
			// Toggle folder open/close
			if f.Cursor < len(f.FlatFiles) {
				node := f.FlatFiles[f.Cursor]
				if node.IsDir {
					node.IsOpen = !node.IsOpen
					// Rebuild the flattened list
					root := f.findRoot()
					f.FlatFiles = file.FlattenFileTree(root)
					f.updateViewportContent()
				}
			}
		case " ":
			// Toggle file selection
			if f.Cursor < len(f.FlatFiles) {
				node := f.FlatFiles[f.Cursor]
				if !node.IsDir {
					node.Selected = !node.Selected
					if node.Selected {
						f.Selected = append(f.Selected, node)
					} else {
						// Remove from selected files
						for i, selectedFile := range f.Selected {
							if selectedFile == node {
								f.Selected = append(f.Selected[:i], f.Selected[i+1:]...)
								break
							}
						}
					}
					f.updateViewportContent()
				}
			}
		}
	}

	// Update viewport
	f.Viewport, cmd = f.Viewport.Update(msg)
	return f, cmd
}

func (f *FileSelect) View() string {
	return RenderLayout(
		f.Title,
		f.Viewport.View(),
		"[↑↓ Navigate] [Enter: Toggle folder] [Space: Select file] [Tab: Next]",
		f.Width,
		f.Height,
	)
}

func (f *FileSelect) updateViewportContent() {
	content := ""
	for i, node := range f.FlatFiles {
		cursor := " "
		if i == f.Cursor {
			cursor = ">"
		}

		line := file.RenderFileNode(node)
		if i == f.Cursor {
			line = selectedStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	selectedInfo := fmt.Sprintf("\nSelected: %d files", len(f.Selected))
	content += selectedInfo

	f.Viewport.SetContent(content)
}

func (f *FileSelect) ensureCursorVisible() {
	// Calculate cursor position in viewport
	lineHeight := 1
	cursorPosition := f.Cursor * lineHeight

	// Scroll to make cursor visible
	if cursorPosition < f.Viewport.YOffset {
		f.Viewport.YOffset = cursorPosition
	} else if cursorPosition >= f.Viewport.YOffset+f.Viewport.Height {
		f.Viewport.YOffset = cursorPosition - f.Viewport.Height + 1
	}
}

func (f *FileSelect) findRoot() *file.FileNode {
	// Find the root node by traversing up from any flat file
	if len(f.FlatFiles) > 0 {
		node := f.FlatFiles[0]
		for node.Parent != nil {
			node = node.Parent
		}
		return node
	}
	return nil
}

func (f *FileSelect) Next() (Component, tea.Cmd) {
	if len(f.Selected) > 0 {
		// Create file template selection component with current dimensions
		templates := []string{"Code Review", "Documentation", "Custom..."}
		return NewTemplateSelect("file", templates, f.Selected, f.Width, f.Height), nil
	}
	return f, nil
}

func (f *FileSelect) Prev() (Component, tea.Cmd) {
	// Return to prompt type with current dimensions
	return NewPromptType(f.Width, f.Height), nil
}

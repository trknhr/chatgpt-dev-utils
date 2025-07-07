package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type Edit struct {
	PromptType       string
	SelectedTemplate string
	TemplateContent  string
	SelectedFiles    []*file.FileNode
	Textarea         textarea.Model
	Width            int
	Height           int
}

func NewEdit(promptType, selectedTemplate, templateContent string, selectedFiles []*file.FileNode, width, height int) *Edit {
	ta := textarea.New()
	ta.Placeholder = "Edit your prompt here..."
	ta.Focus()
	ta.SetWidth(width - 6)
	ta.SetHeight(10)
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Width(0)
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Width(0)
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 1)
	ta.SetValue(templateContent)

	return &Edit{
		PromptType:       promptType,
		SelectedTemplate: selectedTemplate,
		TemplateContent:  templateContent,
		SelectedFiles:    selectedFiles,
		Textarea:         ta,
		Width:            width,
		Height:           height,
	}
}

func (e *Edit) Init() tea.Cmd {
	return textarea.Blink
}

func (e *Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.Width = msg.Width
		e.Height = msg.Height
		// Update textarea dimensions
		textareaHeight := e.Height - 9
		if textareaHeight < 5 {
			textareaHeight = 5
		}
		e.Textarea.SetHeight(textareaHeight)
		boxWidth := e.Width - 4
		if e.Width < 90 {
			boxWidth = e.Width - 10
		}
		e.Textarea.SetWidth(boxWidth - 2)

	case tea.KeyMsg:
		if msg.Type == tea.KeyTab {
			// Move to final step
			return e, nil
		}
	}

	e.Textarea, cmd = e.Textarea.Update(msg)
	return e, cmd
}

func (e *Edit) View() string {
	if e.Height < 10 || e.Width < 20 {
		return "Your terminal is too small."
	}

	var title string
	if e.PromptType == "git" {
		title = "Step 3: Review & Edit"
	} else {
		title = "Step 4: Review & Edit"
	}

	// Update textarea dimensions for current view
	textareaHeight := e.Height - 10
	if textareaHeight < 5 {
		textareaHeight = 5
	}
	e.Textarea.SetHeight(textareaHeight)

	boxWidth := e.Width - 4
	e.Textarea.SetWidth(boxWidth - 2)

	body := e.Textarea.View()

	return RenderLayout(
		title,
		body,
		"[↑↓←→ Type freely] [Tab: Next] [Esc: Back]",
		e.Width,
		e.Height,
	)
}

func (e *Edit) Next() (Component, tea.Cmd) {
	finalPrompt := e.Textarea.Value()
	// Note: WebSocket context will be injected by Root component
	return NewFinal(e.PromptType, e.SelectedTemplate, finalPrompt, e.SelectedFiles, e.Width, e.Height, false, nil, nil), nil
}

func (e *Edit) Prev() (Component, tea.Cmd) {
	// Go back to template selection
	var templates []string
	if e.PromptType == "git" {
		templates = []string{"Code Review", "Commit Message", "Change Summary", "Custom..."}
	} else {
		templates = []string{"Code Review", "Documentation", "Custom..."}
	}
	return NewTemplateSelect(e.PromptType, templates, e.SelectedFiles, e.Width, e.Height), nil
}

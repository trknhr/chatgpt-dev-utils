package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type TemplateSelect struct {
	PromptType    string // "file" or "git"
	Templates     []string
	Cursor        int
	SelectedFiles []*file.FileNode // Only used for file prompts
	Width         int
	Height        int
}

func NewTemplateSelect(promptType string, templates []string, selectedFiles []*file.FileNode, width, height int) *TemplateSelect {
	return &TemplateSelect{
		PromptType:    promptType,
		Templates:     templates,
		Cursor:        0,
		SelectedFiles: selectedFiles,
		Width:         width,
		Height:        height,
	}
}

func (t *TemplateSelect) Init() tea.Cmd {
	return nil
}

func (t *TemplateSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.Width = msg.Width
		t.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.Cursor > 0 {
				t.Cursor--
			}
		case "down", "j":
			if t.Cursor < len(t.Templates)-1 {
				t.Cursor++
			}
		}
	}
	return t, nil
}

func (t *TemplateSelect) View() string {
	var title string
	if t.PromptType == "git" {
		title = "Step 2: Choose Prompt Template"
	} else {
		title = "Step 3: Choose Prompt Template"
	}

	content := ""
	for i, template := range t.Templates {
		cursor := " "
		if t.Cursor == i {
			cursor = ">"
			template = selectedStyle.Render(template)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, template)
	}

	return RenderLayout(
		title,
		content,
		"[↑↓ Navigate] [Tab: Next] [Esc: Back]",
		t.Width,
		t.Height,
	)
}

func (t *TemplateSelect) Next() (Component, tea.Cmd) {
	selectedTemplate := t.Templates[t.Cursor]

	// Get the appropriate template content
	var templateContent string
	if t.PromptType == "git" {
		gitTemplates := map[string]string{
			"Code Review":    "Please review this diff and provide feedback:\n\n$(git diff --cached)\n\nFocus on:\n- Code quality\n- Security issues\n- Performance considerations",
			"Commit Message": "Generate a concise commit message for the following staged changes:\n```\n$(git diff --cached)\n```\n\nFollow the format used in recent commits:\n```\n$(git log -n 3 --pretty=format:%s)\n```\n\nFormat: type(scope): description\n\nOnly return the commit message in plain text. Do not include explanations or comments.",
			"Change Summary": "Summarize the changes in this commit:\n\n$(git log --oneline -1)\n$(git diff HEAD~1)",
			"Custom...":      "$(git diff --cached)",
		}
		templateContent = gitTemplates[selectedTemplate]
	} else {
		fileTemplates := map[string]string{
			"Code Review":   "Please review this code and provide feedback:\n\n$(files)\n\nFocus on:\n- Code quality\n- Best practices\n- Potential issues",
			"Documentation": "Generate documentation for this code:\n\n$(files)\n\nInclude:\n- Function descriptions\n- Usage examples\n- Parameters and return values",
			"Custom...":     "Please add your prompt with $(files)",
		}
		templateContent = fileTemplates[selectedTemplate]
	}

	// Create edit component with WebSocket context placeholder
	return NewEdit(t.PromptType, selectedTemplate, templateContent, t.SelectedFiles, t.Width, t.Height), nil
}

func (t *TemplateSelect) Prev() (Component, tea.Cmd) {
	if t.PromptType == "file" {
		// Go back to file selection
		root := file.BuildFileTree(".")
		flat := file.FlattenFileTree(root)
		vp := viewport.New(t.Width-4, t.Height-8)
		return NewFileSelect(flat, t.SelectedFiles, vp, 0, t.Width, t.Height, ""), nil
	}
	// Go back to prompt type selection
	return NewPromptType(t.Width, t.Height), nil
}

package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

type Step int

const (
	StepPromptType Step = iota
	StepFileSelect
	StepGitTemplate
	StepFileTemplate
	StepGitEdit
	StepFileEdit
	StepFinal
)

type CheckConnectionMsg struct{}

type Model struct {
	CurrentStep        Step
	PromptType         string // "file" or "git"
	FileTree           *file.FileNode
	FlatFiles          []*file.FileNode // flattened view for navigation
	Cursor             int
	SelectedFiles      []*file.FileNode
	Templates          []string
	SelectedTemplate   string
	CustomPrompt       string
	FinalPrompt        string
	GitTemplates       map[string]string
	FileTemplates      map[string]string
	Message            string
	Width              int
	Height             int
	Viewport           viewport.Model
	ExtensionConnected bool // Track if extension is connected via WebSocket
	Textarea           textarea.Model
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(50)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(1, 0)
)

// InitialModel creates the initial TUI model
func InitialModel() Model {
	gitTemplates := map[string]string{
		"Code Review":    "Please review this diff and provide feedback:\n\n$(git diff --cached)\n\nFocus on:\n- Code quality\n- Security issues\n- Performance considerations",
		"Commit Message": "Generate a concise commit message for the following staged changes:\n\n$(git diff --cached)\n\nFollow the format used in recent commits:\n\n$(git log -n 3 --pretty=format:'%h %s')\n\nFormat: type(scope): description\n\nOnly return the commit message in plain text. Do not include explanations or comments.",
		"Change Summary": "Summarize the changes in this commit:\n\n$(git log --oneline -1)\n$(git diff HEAD~1)",
		"Custom...":      "$(git diff --cached)",
	}

	fileTemplates := map[string]string{
		"Code Review":   "Please review this code and provide feedback:\n\n$(files)\n\nFocus on:\n- Code quality\n- Best practices\n- Potential issues",
		"Documentation": "Generate documentation for this code:\n\n$(files)\n\nInclude:\n- Function descriptions\n- Usage examples\n- Parameters and return values",
		"Custom...":     "Please add your prompt with $(files)",
	}

	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	ta := textarea.New()
	ta.Placeholder = "Edit your prompt here..."
	ta.Focus()
	ta.SetWidth(60 - 4)
	ta.SetHeight(10)
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Width(0)
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Width(0)
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252"))
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 1)

	return Model{
		CurrentStep:        StepPromptType,
		GitTemplates:       gitTemplates,
		FileTemplates:      fileTemplates,
		Viewport:           vp,
		ExtensionConnected: false,
		Textarea:           ta,
	}
}

// ... Move all Model methods and step-related functions here ...

// Model struct, initialModel, Init, Update, View, step-related update/view functions, styles, and helpers from main.go
// (Copy all TUI logic here, export as needed)

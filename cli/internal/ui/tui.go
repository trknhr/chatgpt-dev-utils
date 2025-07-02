package ui

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
	"github.com/trknhr/chatgpt-dev-utils/internal/utils"
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
	BroadcastChan      chan<- string
	ClientsCount       func() int
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 0)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 0).
			Width(50)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 0).
			Margin(0, 0)
)

// InitialModel creates the initial TUI model
func InitialModel(broadcastChan chan<- string, clientsCount func() int) Model {
	gitTemplates := map[string]string{
		"Code Review":    "Please review this diff and provide feedback:\n\n$(git diff --cached)\n\nFocus on:\n- Code quality\n- Security issues\n- Performance considerations",
		"Commit Message": "Generate a concise commit message for the following staged changes:\n```\n$(git diff --cached)\n```\n\nFollow the format used in recent commits:\n```\n$(git log -n 3 --pretty=format:%s)\n```\n\nFormat: type(scope): description\n\nOnly return the commit message in plain text. Do not include explanations or comments.",
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
		BroadcastChan:      broadcastChan,
		ClientsCount:       clientsCount,
	}
}

func (m Model) Init() tea.Cmd {
	// Start a timer to check connection status every 2 seconds
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return CheckConnectionMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case CheckConnectionMsg:
		// Check if there are any connected clients
		connected := m.ClientsCount() > 0
		if connected != m.ExtensionConnected {
			m.ExtensionConnected = connected
		}
		// Restart the timer
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return CheckConnectionMsg{}
		})

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m = m.applyWindowSize()
		return m, nil

	case tea.KeyMsg:
		// Handle viewport navigation when in file select
		if m.CurrentStep == StepFileSelect {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				if m.CurrentStep > StepPromptType {
					m.CurrentStep--
					m.Cursor = 0
				}
				return m, nil
			default:
				return m.updateFileSelect(msg)
			}
		}

		// Handle other steps
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.CurrentStep == StepGitTemplate {
				m.CurrentStep = StepPromptType
				m.Cursor = 0
			} else if m.CurrentStep > StepPromptType {
				m.CurrentStep--
				m.Cursor = 0
			}
			return m, nil
		default:
			switch m.CurrentStep {
			case StepPromptType:
				return m.updatePromptType(msg)
			case StepGitTemplate, StepFileTemplate:
				return m.updateTemplateSelect(msg)
			case StepGitEdit:
				return m.updateGitEdit(msg)
			case StepFileEdit:
				return m.updateFileEdit(msg)
			case StepFinal:
				return m.updateFinal(msg)
			}
		}
	}

	// Update viewport
	if m.CurrentStep == StepFileSelect {
		m.Viewport, cmd = m.Viewport.Update(msg)
	}

	return m, cmd
}

func (m Model) applyWindowSize() Model {
	if m.CurrentStep == StepFileSelect {
		headerHeight := 4
		footerHeight := 4
		viewportHeight := m.Height - headerHeight - footerHeight
		if viewportHeight < 3 {
			viewportHeight = 3
		}
		viewportWidth := m.Width - 4
		if viewportWidth < 20 {
			viewportWidth = 20
		}

		m.Viewport.Width = viewportWidth
		m.Viewport.Height = viewportHeight
	}

	return m
}

func (m Model) updatePromptType(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < 1 {
			m.Cursor++
		}
	case "tab":
		if m.Cursor == 0 {
			m.PromptType = "file"
			m.CurrentStep = StepFileSelect
			m.Cursor = 0
			// Build file tree
			m.FileTree = file.BuildFileTree(".")
			m.FlatFiles = file.FlattenFileTree(m.FileTree)
			// Initialize viewport content
			m.updateViewportContent()
			m = m.applyWindowSize()
		} else {
			m.PromptType = "git"
			m.CurrentStep = StepGitTemplate
			m.Cursor = 0
			m.Templates = []string{"Code Review", "Commit Message", "Change Summary", "Custom..."}
		}
	}
	return m, nil
}

func (m Model) updateFileSelect(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			m.updateViewportContent()
			m.ensureCursorVisible()
		}
	case "down", "j":
		if m.Cursor < len(m.FlatFiles)-1 {
			m.Cursor++
			m.updateViewportContent()
			m.ensureCursorVisible()
		}
	case "enter":
		// Toggle folder open/close
		if m.Cursor < len(m.FlatFiles) {
			node := m.FlatFiles[m.Cursor]
			if node.IsDir {
				node.IsOpen = !node.IsOpen
				m.FlatFiles = file.FlattenFileTree(m.FileTree)
				m.updateViewportContent()
			}
		}
	case " ":
		// Toggle file selection
		if m.Cursor < len(m.FlatFiles) {
			node := m.FlatFiles[m.Cursor]
			if !node.IsDir {
				node.Selected = !node.Selected
				if node.Selected {
					m.SelectedFiles = append(m.SelectedFiles, node)
				} else {
					// Remove from selected files
					for i, f := range m.SelectedFiles {
						if f == node {
							m.SelectedFiles = append(m.SelectedFiles[:i], m.SelectedFiles[i+1:]...)
							break
						}
					}
				}
				m.updateViewportContent()
			}
		}
	case "tab":
		if len(m.SelectedFiles) > 0 {
			m.CurrentStep = StepFileTemplate
			m.Cursor = 0
			m.Templates = []string{"Code Review", "Documentation", "Custom..."}
		}
		m = m.applyWindowSize()
	}
	return m, nil
}

func (m *Model) updateViewportContent() {
	content := ""
	for i, node := range m.FlatFiles {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}

		line := file.RenderFileNode(node)
		if i == m.Cursor {
			line = selectedStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	selectedInfo := fmt.Sprintf("\nSelected: %d files", len(m.SelectedFiles))
	content += selectedInfo

	m.Viewport.SetContent(content)
}

func (m *Model) ensureCursorVisible() {
	// Calculate cursor position in viewport
	lineHeight := 1
	cursorPosition := m.Cursor * lineHeight

	// Scroll to make cursor visible
	if cursorPosition < m.Viewport.YOffset {
		m.Viewport.YOffset = cursorPosition
	} else if cursorPosition >= m.Viewport.YOffset+m.Viewport.Height {
		m.Viewport.YOffset = cursorPosition - m.Viewport.Height + 1
	}
}

func (m Model) updateTemplateSelect(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Templates)-1 {
			m.Cursor++
		}
	case "tab":
		m.SelectedTemplate = m.Templates[m.Cursor]
		switch m.PromptType {
		case "git":
			m.CurrentStep = StepGitEdit
			m.CustomPrompt = m.GitTemplates[m.SelectedTemplate]
			m.Textarea.SetValue(m.CustomPrompt)
			m.Textarea.Focus()
		case "file":
			m.CurrentStep = StepFileEdit
			m.CustomPrompt = m.FileTemplates[m.SelectedTemplate]
			m.Textarea.SetValue(m.CustomPrompt)
			m.Textarea.Focus()
		}
		m.Cursor = 0
		m = m.applyWindowSize()
	}
	return m, nil
}

func (m Model) updateGitEdit(msg tea.KeyMsg) (Model, tea.Cmd) {
	if msg.Type == tea.KeyTab {
		m.CustomPrompt = m.Textarea.Value()
		m.FinalPrompt = m.CustomPrompt
		m.CurrentStep = StepFinal
		return m, nil
	}

	var cmd tea.Cmd
	m.Textarea, cmd = m.Textarea.Update(msg)
	return m, cmd
}

func (m Model) updateFileEdit(msg tea.KeyMsg) (Model, tea.Cmd) {
	if msg.Type == tea.KeyTab {
		m.CustomPrompt = m.Textarea.Value()
		m.FinalPrompt = m.CustomPrompt
		m.CurrentStep = StepFinal
		return m, nil
	}

	var cmd tea.Cmd
	m.Textarea, cmd = m.Textarea.Update(msg)
	return m, cmd
}

func (m Model) updateFinal(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		var finalContent string
		if m.PromptType == "file" {
			// Generate file prompt with actual content
			finalContent = utils.GenerateFilePrompt(m.FinalPrompt, m.SelectedFiles)
		} else {
			// Execute git commands
			finalContent = utils.ExecuteGitCommands(m.FinalPrompt)
		}
		clipboard.WriteAll(finalContent)
		m.Message = "Copied to clipboard!"
	case "e":
		if m.ExtensionConnected {
			var finalContent string
			if m.PromptType == "file" {
				// Generate file prompt with actual content
				finalContent = utils.GenerateFilePrompt(m.FinalPrompt, m.SelectedFiles)
			} else {
				// Execute git commands
				finalContent = utils.ExecuteGitCommands(m.FinalPrompt)
			}

			payload := map[string]string{
				"type":   "chatgpt-prompt",
				"prompt": finalContent,
			}

			jsonBytes, err := json.Marshal(payload)

			if err != nil {
				return m, nil
			}
			// Send to extension via WebSocket
			select {
			case m.BroadcastChan <- string(jsonBytes):
				m.Message = "Sent to extension!"
			default:
				m.Message = "Extension not connected"
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.CurrentStep {
	case StepPromptType:
		return m.viewPromptType()
	case StepFileSelect:
		return m.viewFileSelect()
	case StepGitTemplate, StepFileTemplate:
		return m.viewTemplateSelect()
	case StepGitEdit:
		return m.viewGitEdit()
	case StepFileEdit:
		return m.viewFileEdit()
	case StepFinal:
		return m.viewFinal()
	}
	return ""
}

func (m Model) viewPromptType() string {
	title := titleStyle.Render("Step 1/4: Choose Prompt Type")

	options := []string{
		"File based Prompt",
		"Git based Prompt",
	}

	content := ""
	for i, option := range options {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
			option = selectedStyle.Render(option)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, option)
	}

	help := helpStyle.Render("[↑↓ Navigate] [Tab: Next] [Ctrl+C: Quit]")

	// Adjust box width based on terminal size
	boxWidth := m.Width - 4
	if m.Width > 0 && m.Width < 60 {
		boxWidth = m.Width - 10
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		boxStyle.Width(boxWidth).Render(content),
		help,
		m.Message,
	)
}

func (m Model) viewFileSelect() string {
	title := titleStyle.Render("Step 2/4: Select Files")
	help := helpStyle.Render("[↑↓ Navigate] [Enter: Toggle folder] [Space: Select file] [Tab: Next]")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		m.Viewport.View(),
		help,
		m.Message,
	)
}

func (m Model) viewTemplateSelect() string {
	var title string
	if m.PromptType == "git" {
		title = titleStyle.Render("Step 2/4: Choose Prompt Template")
	} else {
		title = titleStyle.Render("Step 3/4: Choose Prompt Template")
	}

	content := ""
	for i, template := range m.Templates {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
			template = selectedStyle.Render(template)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, template)
	}

	help := helpStyle.Render("[↑↓ Navigate] [Tab: Next] [Esc: Back]")

	// Adjust box width based on terminal size
	boxWidth := m.Width - 4
	if m.Width > 0 && m.Width < 60 {
		boxWidth = m.Width - 10
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		boxStyle.Width(boxWidth).Render(content),
		help,
		m.Message,
	)
}

func (m Model) viewGitEdit() string {
	if m.Height < 10 || m.Width < 20 {
		return "Your terminal is too small."
	}

	title := titleStyle.Render("Step 3/4: Review & Edit")
	help := helpStyle.Render("[↑↓←→ Type freely] [Tab: Next] [Esc: Back]")

	textareaHeight := m.Height - 9
	if textareaHeight < 5 {
		textareaHeight = 5
	}
	m.Textarea.SetHeight(textareaHeight)

	boxWidth := m.Width - 4
	if m.Width < 90 {
		boxWidth = m.Width - 10
	}
	m.Textarea.SetWidth(boxWidth - 2)

	body := fmt.Sprintf("Template: %s\n%s", m.SelectedTemplate, m.Textarea.View())
	boxContent := boxStyle.Width(boxWidth).Render(body)

	content := []string{
		title,
		boxContent,
		help,
	}
	if m.Message != "" {
		content = append(content, m.Message)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) viewFileEdit() string {
	if m.Height < 10 || m.Width < 20 {
		return "Your terminal is too small."
	}

	title := titleStyle.Render("Step 3/4: Review & Edit")
	help := helpStyle.Render("[↑↓←→ Type freely] [Tab: Next] [Esc: Back]")

	textareaHeight := m.Height - 9

	if textareaHeight < 5 {
		textareaHeight = 5
	}
	m.Textarea.SetHeight(textareaHeight)

	boxWidth := m.Width - 4

	// Subtract 2 for boxStyle's left/right padding or border
	m.Textarea.SetWidth(boxWidth - 2)
	if m.Width < 90 {
		boxWidth = m.Width - 10
	}

	body := fmt.Sprintf("Template: %s\n%s", m.SelectedTemplate, m.Textarea.View())
	boxContent := boxStyle.Width(boxWidth).Render(body)

	content := []string{
		title,
		boxContent,
		help,
	}
	if m.Message != "" {
		content = append(content, m.Message)
	}

	contentView := lipgloss.JoinVertical(lipgloss.Left, content...)
	return lipgloss.NewStyle().Margin(0, 0).Render(contentView)
}

func (m Model) viewFinal() string {
	var title string
	if m.PromptType == "file" {
		title = titleStyle.Render("Step 4/4: Review & Copy")
	} else {
		title = titleStyle.Render("Step 4/4: Copy Prompt")
	}

	var content string
	if m.PromptType == "file" {
		// Show template with selected files list
		template := m.FinalPrompt
		if template == "" {
			template = "Please analyze these files:\n\n$(files)"
		}

		// Build selected files list
		filesList := "Selected files:\n"
		for _, file := range m.SelectedFiles {
			filesList += fmt.Sprintf("- %s\n", file.Path)
		}

		content = fmt.Sprintf("Template: %s\n\n%s\n\n%s",
			m.SelectedTemplate,
			template,
			filesList)
	} else {
		// Git-based: show the prompt as before
		preview := m.FinalPrompt
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		content = fmt.Sprintf("Ready to copy:\n\n%s", preview)
	}

	helpStr := "[C: Copy with Content] [Esc: Back]"
	if m.ExtensionConnected {
		helpStr += " [E: Send to Extension]"

	}
	help := helpStyle.Render(helpStr)

	boxWidth := m.Width - 4

	if m.Width > 0 && m.Width < 90 {
		boxWidth = m.Width - 10
	}

	// Don't force height, let content determine it
	boxContent := boxStyle.Height(m.Height - 50).Width(boxWidth).Render(content)

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		title,
		boxContent,
		help,
		m.Message,
	)
}

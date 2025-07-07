package components

import (
	"encoding/json"
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trknhr/chatgpt-dev-utils/internal/file"
	"github.com/trknhr/chatgpt-dev-utils/internal/utils"
)

type Final struct {
	PromptType         string
	SelectedTemplate   string
	FinalPrompt        string
	SelectedFiles      []*file.FileNode
	Width              int
	Height             int
	Message            string
	ExtensionConnected bool
	BroadcastChan      chan<- string
	ClientsCount       func() int
}

func NewFinal(promptType, selectedTemplate, finalPrompt string, selectedFiles []*file.FileNode, width, height int, extensionConnected bool, broadcastChan chan<- string, clientsCount func() int) *Final {
	return &Final{
		PromptType:         promptType,
		SelectedTemplate:   selectedTemplate,
		FinalPrompt:        finalPrompt,
		SelectedFiles:      selectedFiles,
		Width:              width,
		Height:             height,
		ExtensionConnected: extensionConnected,
		BroadcastChan:      broadcastChan,
		ClientsCount:       clientsCount,
	}
}

func (f *Final) Init() tea.Cmd {
	return nil
}

func (f *Final) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.Width = msg.Width
		f.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			var finalContent string
			if f.PromptType == "file" {
				// Generate file prompt with actual content
				finalContent = utils.GenerateFilePrompt(f.FinalPrompt, f.SelectedFiles)
			} else {
				// Execute git commands
				finalContent = utils.ExecuteGitCommands(f.FinalPrompt)
			}
			clipboard.WriteAll(finalContent)
			f.Message = "Copied to clipboard!"
		case "e":
			if f.ExtensionConnected && f.BroadcastChan != nil {
				var finalContent string
				if f.PromptType == "file" {
					// Generate file prompt with actual content
					finalContent = utils.GenerateFilePrompt(f.FinalPrompt, f.SelectedFiles)
				} else {
					// Execute git commands
					finalContent = utils.ExecuteGitCommands(f.FinalPrompt)
				}

				payload := map[string]string{
					"type":   "chatgpt-prompt",
					"prompt": finalContent,
				}

				jsonBytes, err := json.Marshal(payload)

				if err != nil {
					f.Message = "Error marshaling JSON"
					return f, nil
				}
				// Send to extension via WebSocket
				select {
				case f.BroadcastChan <- string(jsonBytes):
					f.Message = "Sent to extension!"
				default:
					f.Message = "Extension not connected"
				}
			}
		}
	case CheckConnectionMsg:
		// Update extension connection status
		if f.ClientsCount != nil {
			f.ExtensionConnected = f.ClientsCount() > 0
		}
	}
	return f, nil
}

func (f *Final) View() string {
	var title string
	if f.PromptType == "file" {
		title = "Step 5: Copy Prompt"
	} else {
		title = "Step 4: Copy Prompt"
	}

	var content string
	if f.PromptType == "file" {
		// Show template with selected files list
		template := f.FinalPrompt
		if template == "" {
			template = "Please analyze these files:\n\n$(files)"
		}

		// Build selected files list
		filesList := "Selected files:\n"
		for _, file := range f.SelectedFiles {
			filesList += fmt.Sprintf("- %s\n", file.Path)
		}

		content = fmt.Sprintf("Template: %s\n\n%s\n\n%s",
			f.SelectedTemplate,
			template,
			filesList)
	} else {
		// Git-based: show the prompt as before
		preview := f.FinalPrompt
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		content = fmt.Sprintf("Ready to copy:\n\n%s", preview)
	}

	helpStr := "[C: Copy with Content] [Esc: Back]"
	if f.ExtensionConnected {
		helpStr += " [E: Send to Extension]"
	}

	return RenderLayoutWithMessage(
		title,
		content,
		helpStr,
		f.Message,
		f.Width,
		f.Height,
	)
}

func (f *Final) Next() (Component, tea.Cmd) {
	// Final step, no next
	return f, nil
}

func (f *Final) Prev() (Component, tea.Cmd) {
	// Go back to edit step
	var templateContent string
	if f.PromptType == "git" {
		gitTemplates := map[string]string{
			"Code Review":    "Please review this diff and provide feedback:\n\n$(git diff --cached)\n\nFocus on:\n- Code quality\n- Security issues\n- Performance considerations",
			"Commit Message": "Generate a concise commit message for the following staged changes:\n```\n$(git diff --cached)\n```\n\nFollow the format used in recent commits:\n```\n$(git log -n 3 --pretty=format:%s)\n```\n\nFormat: type(scope): description\n\nOnly return the commit message in plain text. Do not include explanations or comments.",
			"Change Summary": "Summarize the changes in this commit:\n\n$(git log --oneline -1)\n$(git diff HEAD~1)",
			"Custom...":      "$(git diff --cached)",
		}
		templateContent = gitTemplates[f.SelectedTemplate]
	} else {
		fileTemplates := map[string]string{
			"Code Review":   "Please review this code and provide feedback:\n\n$(files)\n\nFocus on:\n- Code quality\n- Best practices\n- Potential issues",
			"Documentation": "Generate documentation for this code:\n\n$(files)\n\nInclude:\n- Function descriptions\n- Usage examples\n- Parameters and return values",
			"Custom...":     "Please add your prompt with $(files)",
		}
		templateContent = fileTemplates[f.SelectedTemplate]
	}

	return NewEdit(f.PromptType, f.SelectedTemplate, templateContent, f.SelectedFiles, f.Width, f.Height), nil
}

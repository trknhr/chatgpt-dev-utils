package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

// WebSocket related variables
var (
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
)

// App states
type step int

const (
	stepPromptType step = iota
	stepFileSelect
	stepGitTemplate
	stepFileTemplate
	stepGitEdit
	stepFinal
)

// File tree node
type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	IsOpen   bool
	Selected bool
	Children []*FileNode
	Parent   *FileNode
}

// Model represents the application state
type Model struct {
	currentStep      step
	promptType       string // "file" or "git"
	fileTree         *FileNode
	flatFiles        []*FileNode // flattened view for navigation
	cursor           int
	selectedFiles    []*FileNode
	templates        []string
	selectedTemplate string
	customPrompt     string
	finalPrompt      string
	gitTemplates     map[string]string
	fileTemplates    map[string]string
	message          string
	width            int
	height           int
	viewport         viewport.Model
	extensionEnabled bool // Flag to show extension button
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

func initialModel(extensionEnabled bool) Model {
	// Initialize templates
	gitTemplates := map[string]string{
		"Code Review":    "Please review this diff and provide feedback:\n\n$(git diff --cached)\n\nFocus on:\n- Code quality\n- Security issues\n- Performance considerations",
		"Commit Message": "Generate a concise commit message for these changes:\n\n$(git diff --cached)\n\nFormat: type(scope): description",
		"Change Summary": "Summarize the changes in this commit:\n\n$(git log --oneline -1)\n$(git diff HEAD~1)",
		"Custom...":      "",
	}

	fileTemplates := map[string]string{
		"Code Review":   "Please review this code and provide feedback:\n\n$(files)\n\nFocus on:\n- Code quality\n- Best practices\n- Potential issues",
		"Documentation": "Generate documentation for this code:\n\n$(files)\n\nInclude:\n- Function descriptions\n- Usage examples\n- Parameters and return values",
		"Custom...":     "",
	}

	// Initialize viewport
	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return Model{
		currentStep:      stepPromptType,
		gitTemplates:     gitTemplates,
		fileTemplates:    fileTemplates,
		viewport:         vp,
		extensionEnabled: extensionEnabled,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport size when on file select step
		if m.currentStep == stepFileSelect {
			headerHeight := 4 // Title + spacing
			footerHeight := 4 // Help + message + spacing
			viewportHeight := msg.Height - headerHeight - footerHeight
			if viewportHeight < 3 {
				viewportHeight = 3
			}

			viewportWidth := msg.Width - 4 // Account for border and padding
			if viewportWidth < 20 {
				viewportWidth = 20
			}

			m.viewport.Width = viewportWidth
			m.viewport.Height = viewportHeight
		}
		return m, nil

	case tea.KeyMsg:
		// Handle viewport navigation when in file select
		if m.currentStep == stepFileSelect {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				if m.currentStep > stepPromptType {
					m.currentStep--
					m.cursor = 0
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
			if m.currentStep > stepPromptType {
				m.currentStep--
				m.cursor = 0
			}
			return m, nil
		}

		switch m.currentStep {
		case stepPromptType:
			return m.updatePromptType(msg)
		case stepGitTemplate, stepFileTemplate:
			return m.updateTemplateSelect(msg)
		case stepGitEdit:
			return m.updateGitEdit(msg)
		case stepFinal:
			return m.updateFinal(msg)
		}
	}

	// Update viewport
	if m.currentStep == stepFileSelect {
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m Model) updatePromptType(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 1 {
			m.cursor++
		}
	case "enter", " ":
		if m.cursor == 0 {
			m.promptType = "file"
			m.currentStep = stepFileSelect
			m.cursor = 0
			// Build file tree
			m.fileTree = m.buildFileTree(".")
			m.flatFiles = m.flattenFileTree(m.fileTree)
			// Initialize viewport content
			m.updateViewportContent()
		} else {
			m.promptType = "git"
			m.currentStep = stepGitTemplate
			m.cursor = 0
			m.templates = []string{"Code Review", "Commit Message", "Change Summary", "Custom..."}
		}
	}
	return m, nil
}

func (m Model) updateFileSelect(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.updateViewportContent()
			m.ensureCursorVisible()
		}
	case "down", "j":
		if m.cursor < len(m.flatFiles)-1 {
			m.cursor++
			m.updateViewportContent()
			m.ensureCursorVisible()
		}
	case "enter":
		// Toggle folder open/close
		if m.cursor < len(m.flatFiles) {
			node := m.flatFiles[m.cursor]
			if node.IsDir {
				node.IsOpen = !node.IsOpen
				m.flatFiles = m.flattenFileTree(m.fileTree)
				m.updateViewportContent()
			}
		}
	case " ":
		// Toggle file selection
		if m.cursor < len(m.flatFiles) {
			node := m.flatFiles[m.cursor]
			if !node.IsDir {
				node.Selected = !node.Selected
				if node.Selected {
					m.selectedFiles = append(m.selectedFiles, node)
				} else {
					// Remove from selected files
					for i, f := range m.selectedFiles {
						if f == node {
							m.selectedFiles = append(m.selectedFiles[:i], m.selectedFiles[i+1:]...)
							break
						}
					}
				}
				m.updateViewportContent()
			}
		}
	case "tab":
		if len(m.selectedFiles) > 0 {
			m.currentStep = stepFileTemplate
			m.cursor = 0
			m.templates = []string{"Code Review", "Documentation", "Custom..."}
		}
	}
	return m, nil
}

func (m *Model) updateViewportContent() {
	content := ""
	for i, node := range m.flatFiles {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		line := m.renderFileNode(node)
		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		content += fmt.Sprintf("%s %s\n", cursor, line)
	}

	selectedInfo := fmt.Sprintf("\nSelected: %d files", len(m.selectedFiles))
	content += selectedInfo

	m.viewport.SetContent(content)
}

func (m *Model) ensureCursorVisible() {
	// Calculate cursor position in viewport
	lineHeight := 1
	cursorPosition := m.cursor * lineHeight

	// Scroll to make cursor visible
	if cursorPosition < m.viewport.YOffset {
		m.viewport.YOffset = cursorPosition
	} else if cursorPosition >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.YOffset = cursorPosition - m.viewport.Height + 1
	}
}

func (m Model) updateTemplateSelect(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.templates)-1 {
			m.cursor++
		}
	case "enter", " ":
		m.selectedTemplate = m.templates[m.cursor]
		if m.promptType == "git" {
			m.currentStep = stepGitEdit
			m.customPrompt = m.gitTemplates[m.selectedTemplate]
		} else {
			m.currentStep = stepFinal
			m.finalPrompt = m.generateFilePrompt()
		}
		m.cursor = 0
	}
	return m, nil
}

func (m Model) updateGitEdit(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.currentStep = stepFinal
		m.finalPrompt = m.customPrompt
	}
	return m, nil
}

func (m Model) updateFinal(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		var finalContent string
		if m.promptType == "file" {
			// Generate file prompt with actual content
			finalContent = m.generateFilePrompt()
		} else {
			// Execute git commands
			finalContent = m.executeGitCommands(m.finalPrompt)
		}
		clipboard.WriteAll(finalContent)
		m.message = "Copied to clipboard!"
	case "e":
		if m.extensionEnabled {
			var finalContent string
			if m.promptType == "file" {
				// Generate file prompt with actual content
				finalContent = m.generateFilePrompt()
			} else {
				// Execute git commands
				finalContent = m.executeGitCommands(m.finalPrompt)
			}

			payload := map[string]string{
				"type":   "chatgpt-prompt",
				"prompt": finalContent, //out.String(),
			}

			jsonBytes, err := json.Marshal(payload)

			if err != nil {
				return m, nil
			}
			// Send to extension via WebSocket
			select {
			case broadcast <- string(jsonBytes):
				m.message = "Sent to extension!"
			default:
				m.message = "Extension not connected"
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.currentStep {
	case stepPromptType:
		return m.viewPromptType()
	case stepFileSelect:
		return m.viewFileSelect()
	case stepGitTemplate, stepFileTemplate:
		return m.viewTemplateSelect()
	case stepGitEdit:
		return m.viewGitEdit()
	case stepFinal:
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
		if m.cursor == i {
			cursor = ">"
			option = selectedStyle.Render(option)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, option)
	}

	help := helpStyle.Render("[↑↓ Navigate] [Enter: Select] [Ctrl+C: Quit]")

	// Adjust box width based on terminal size
	boxWidth := 50
	if m.width > 0 && m.width < 60 {
		boxWidth = m.width - 10
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		boxStyle.Width(boxWidth).Render(content),
		help,
		m.message,
	)
}

func (m Model) viewFileSelect() string {
	title := titleStyle.Render("Step 2/4: Select Files")
	help := helpStyle.Render("[↑↓ Navigate] [Enter: Toggle folder] [Space: Select file] [Tab: Next]")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		m.viewport.View(),
		help,
		m.message,
	)
}

func (m Model) viewTemplateSelect() string {
	var title string
	if m.promptType == "git" {
		title = titleStyle.Render("Step 2/4: Choose Prompt Template")
	} else {
		title = titleStyle.Render("Step 3/4: Choose Prompt Template")
	}

	content := ""
	for i, template := range m.templates {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			template = selectedStyle.Render(template)
		}
		content += fmt.Sprintf("%s ◯ %s\n", cursor, template)
	}

	help := helpStyle.Render("[↑↓ Navigate] [Enter: Select] [Esc: Back]")

	// Adjust box width based on terminal size
	boxWidth := 50
	if m.width > 0 && m.width < 60 {
		boxWidth = m.width - 10
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		title,
		boxStyle.Width(boxWidth).Render(content),
		help,
		m.message,
	)
}

func (m Model) viewGitEdit() string {
	title := titleStyle.Render("Step 3/4: Review & Edit")

	content := fmt.Sprintf("Template: %s\n\n%s", m.selectedTemplate, m.customPrompt)

	help := helpStyle.Render("[Tab: Next] [Esc: Back] (Edit functionality would be added here)")

	// Adjust box width based on terminal size
	boxWidth := 80
	if m.width > 0 && m.width < 90 {
		boxWidth = m.width - 10
	}

	boxContent := boxStyle.Width(boxWidth).Render(content)

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		title,
		boxContent,
		help,
		m.message,
	)
}

func (m Model) viewFinal() string {
	var title string
	if m.promptType == "file" {
		title = titleStyle.Render("Step 4/4: Review & Copy")
	} else {
		title = titleStyle.Render("Step 4/4: Copy Prompt")
	}

	var content string
	if m.promptType == "file" {
		// Show template with selected files list
		template := m.fileTemplates[m.selectedTemplate]
		if template == "" {
			template = "Please analyze these files:\n\n$(files)"
		}

		// Build selected files list
		filesList := "Selected files:\n"
		for _, file := range m.selectedFiles {
			filesList += fmt.Sprintf("- %s\n", file.Path)
		}

		content = fmt.Sprintf("Template: %s\n\n%s\n\n%s",
			m.selectedTemplate,
			template,
			filesList)
	} else {
		// Git-based: show the prompt as before
		preview := m.finalPrompt
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		content = fmt.Sprintf("Ready to copy:\n\n%s", preview)
	}

	helpStr := "[C: Copy with Content] [Esc: Back]"
	if m.extensionEnabled {
		helpStr += " [E: Send to Extension]"

	}
	help := helpStyle.Render(helpStr)
	// Adjust dimensions based on terminal size
	boxWidth := 80
	if m.width > 0 && m.width < 90 {
		boxWidth = m.width - 10
	}

	// Don't force height, let content determine it
	boxContent := boxStyle.Width(boxWidth).Render(content)

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		title,
		boxContent,
		help,
		m.message,
	)
}

func (m Model) buildFileTree(root string) *FileNode {
	rootNode := &FileNode{
		Name:   filepath.Base(root),
		Path:   root,
		IsDir:  true,
		IsOpen: true,
	}

	m.buildFileTreeRecursive(rootNode, root, 0)
	return rootNode
}

func (m Model) buildFileTreeRecursive(parent *FileNode, path string, depth int) {
	if depth > 3 { // Limit depth to avoid too deep trees
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	// Sort entries: directories first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// Skip hidden files and common ignore patterns
		if strings.HasPrefix(entry.Name(), ".") ||
			entry.Name() == "node_modules" ||
			entry.Name() == "vendor" {
			continue
		}

		childPath := filepath.Join(path, entry.Name())
		child := &FileNode{
			Name:   entry.Name(),
			Path:   childPath,
			IsDir:  entry.IsDir(),
			IsOpen: false,
			Parent: parent,
		}

		parent.Children = append(parent.Children, child)

		if entry.IsDir() {
			m.buildFileTreeRecursive(child, childPath, depth+1)
		}
	}
}

func (m Model) flattenFileTree(root *FileNode) []*FileNode {
	var result []*FileNode
	m.flattenFileTreeRecursive(root, &result, 0)
	return result
}

func (m Model) flattenFileTreeRecursive(node *FileNode, result *[]*FileNode, depth int) {
	if depth > 0 { // Skip root node
		*result = append(*result, node)
	}

	if node.IsDir && node.IsOpen {
		for _, child := range node.Children {
			m.flattenFileTreeRecursive(child, result, depth+1)
		}
	}
}

func (m Model) renderFileNode(node *FileNode) string {
	depth := m.getNodeDepth(node)
	indent := strings.Repeat("  ", depth-1)

	if node.IsDir {
		icon := "▶"
		if node.IsOpen {
			icon = "▼"
		}
		fileCount := ""
		if !node.IsOpen && len(node.Children) > 0 {
			fileCount = fmt.Sprintf(" (%d items)", len(node.Children))
		}
		return fmt.Sprintf("%s%s %s/%s", indent, icon, node.Name, fileCount)
	} else {
		checkbox := "◯"
		if node.Selected {
			checkbox = "◉"
		}
		return fmt.Sprintf("%s  %s %s", indent, checkbox, node.Name)
	}
}

func (m Model) getNodeDepth(node *FileNode) int {
	depth := 0
	current := node
	for current.Parent != nil {
		depth++
		current = current.Parent
	}
	return depth
}

func (m Model) generateFilePrompt() string {
	template := m.fileTemplates[m.selectedTemplate]
	if template == "" {
		template = "Please analyze these files:\n\n$(files)"
	}

	// Replace $(files) with actual file contents
	fileContents := ""
	for _, file := range m.selectedFiles {
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fileContents += fmt.Sprintf("// Error reading %s: %v\n\n", file.Path, err)
			continue
		}
		fileContents += fmt.Sprintf("// File: %s\n%s\n\n", file.Path, string(content))
	}

	return strings.ReplaceAll(template, "$(files)", fileContents)
}

func (m Model) executeGitCommands(prompt string) string {
	// Simple regex-like replacement for $(git ...)
	lines := strings.Split(prompt, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(git ") {
			start := strings.Index(line, "$(git ")
			end := strings.Index(line[start:], ")")
			if end != -1 {
				end += start
				command := line[start+2 : end] // Remove $( and )
				output := m.executeCommand(command)
				lines[i] = strings.Replace(line, line[start:end+1], output, 1)
			}
		}
	}

	return strings.Join(lines, "\n")
}

func (m Model) generateScrollIndicator(startIndex, displayedCount, totalItems, maxDisplayHeight int) string {
	// This function is no longer needed as viewport handles scrolling
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) executeCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing %s: %v", command, err)
	}

	return strings.TrimSpace(string(output))
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("WebSocket server started at :32123")
	go func() {
		err := http.ListenAndServe(":32123", nil)
		if err != nil {
			log.Printf("WebSocket server error: %v", err)
		}
	}()
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// CORS対応など必要ならここに追記
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// log.Println("Extension disconnected")
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Parse command line flags
	extensionEnabled := false
	for _, arg := range os.Args[1:] {
		if arg == "--extension" || arg == "-e" {
			extensionEnabled = true
			break
		}
	}

	// Start WebSocket server if extension is enabled
	if extensionEnabled {
		go startWebSocketServer()
	}

	p := tea.NewProgram(
		initialModel(extensionEnabled),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

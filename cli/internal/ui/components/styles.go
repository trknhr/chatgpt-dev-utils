package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Shared styles used across components
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 0)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 0).
			Margin(0, 0)
)

// RenderLayout renders a standard layout with title, content box, and help text
func RenderLayout(title, content, help string, width, height int) string {
	// Adjust box width based on terminal size
	boxWidth := width - 4
	if width > 0 && width < 60 {
		boxWidth = width - 10
	}

	dynamicBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 0).
		Width(boxWidth)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		dynamicBoxStyle.Render(content),
		helpStyle.Render(help),
	)
}

// RenderLayoutWithMessage adds an optional message line
func RenderLayoutWithMessage(title, content, help, message string, width, height int) string {
	base := RenderLayout(title, content, help, width, height)
	if message != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			base,
			message)
		// fmt.Sprintf("%s\n%s", base, message)
	}
	return base
}

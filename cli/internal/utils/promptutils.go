package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/trknhr/chatgpt-dev-utils/internal/file"
)

// GenerateFilePrompt replaces $(files) in the template with the contents of selected files
func GenerateFilePrompt(text string, selectedFiles []*file.FileNode) string {
	if text == "" {
		text = "Please analyze these files:\n\n$(files)"
	}

	fileContents := ""
	for _, file := range selectedFiles {
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fileContents += fmt.Sprintf("// Error reading %s: %v\n\n", file.Path, err)
			continue
		}
		fileContents += fmt.Sprintf("// File: %s\n%s\n\n", file.Path, string(content))
	}

	return strings.ReplaceAll(text, "$(files)", fileContents)
}

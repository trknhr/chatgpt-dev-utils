package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand runs a shell command and returns its output or error
func ExecuteCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
	if err != nil {
		return fmt.Sprintf("Error executing %s: %s", command, string(output))
	}

	return strings.TrimSpace(string(output))
}

// ExecuteGitCommands replaces $(git ...) in the prompt with the output of the git command
func ExecuteGitCommands(prompt string) string {
	lines := strings.Split(prompt, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(git ") {
			start := strings.Index(line, "$(git ")
			end := strings.Index(line[start:], ")")
			if end != -1 {
				end += start
				command := line[start+2 : end] // Remove $( and )
				output := ExecuteCommand(command)
				// Handle common git errors
				if strings.Contains(output, "not a git repository") ||
					strings.Contains(output, "does not have any commits yet") ||
					strings.Contains(output, "fatal:") {
					output = "[git error: " + output + "]"
				}
				lines[i] = strings.Replace(line, line[start:end+1], output, 1)
			}
		}
	}
	return strings.Join(lines, "\n")
}

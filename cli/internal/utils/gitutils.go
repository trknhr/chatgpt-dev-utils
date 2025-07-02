package utils

import (
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

// ExecuteCommand runs a shell command and returns its output or error
func executeGitCommand(command string) string {
	parser := shellwords.NewParser()
	fields, err := parser.Parse(command)
	if err != nil || len(fields) == 0 || fields[0] != "git" {
		return "[only 'git' commands are allowed]"
	}

	cmd := exec.Command(fields[0], fields[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "[error executing git command]"
	}

	return strings.TrimSpace(string(out))
}

// ExecuteGitCommands replaces $(git ...) in the prompt with the output of the git command
func ExecuteGitCommands(prompt string) string {
	// lines := strings.Split(prompt, "\n")
	// for i, line := range lines {
	// 	if strings.Contains(line, "$(git ") {
	// 		start := strings.Index(line, "$(git ")
	// 		end := strings.Index(line[start:], ")")
	// 		if end != -1 {
	// 			end += start
	// 			command := line[start+2 : end] // Remove $( and )
	// 			output := ExecuteCommand(command)
	// 			// Handle common git errors
	// 			if strings.Contains(output, "not a git repository") ||
	// 				strings.Contains(output, "does not have any commits yet") ||
	// 				strings.Contains(output, "fatal:") {
	// 				output = "[git error: " + output + "]"
	// 			}
	// 			lines[i] = strings.Replace(line, line[start:end+1], output, 1)
	// 		}
	// 	}
	// }
	// return strings.Join(lines, "\n")
	lines := strings.Split(prompt, "\n")
	for i, line := range lines {
		for {
			start := strings.Index(line, "$(")
			if start == -1 {
				break
			}
			end := strings.Index(line[start:], ")")
			if end == -1 {
				break
			}
			end += start
			command := line[start+2 : end]

			output := executeGitCommand(command)
			line = strings.Replace(line, line[start:end+1], output, 1)
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

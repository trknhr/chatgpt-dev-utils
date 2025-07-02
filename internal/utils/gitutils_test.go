package utils

import (
	"strings"
	"testing"
)

func TestExecuteCommand_Success(t *testing.T) {
	output := ExecuteCommand("echo hello")
	if strings.TrimSpace(output) != "hello" {
		t.Errorf("expected 'hello', got '%s'", output)
	}
}

func TestExecuteCommand_Error(t *testing.T) {
	output := ExecuteCommand("git not-a-real-command")
	if !strings.Contains(output, "Error executing") {
		t.Errorf("expected error message, got '%s'", output)
	}
}

func TestExecuteGitCommands_ReplaceGit(t *testing.T) {
	prompt := "This is a test: $(git --version)"
	output := ExecuteGitCommands(prompt)
	if !strings.Contains(output, "git version") && !strings.Contains(output, "[git error:") {
		t.Errorf("expected git version or error, got '%s'", output)
	}
}

func TestExecuteGitCommands_GitError(t *testing.T) {
	prompt := "$(git not-a-real-command)"
	output := ExecuteGitCommands(prompt)
	if !strings.Contains(output, "[git error:") {
		t.Errorf("expected git error, got '%s'", output)
	}
}

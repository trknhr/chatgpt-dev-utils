package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteGitCommands(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		validate func(t *testing.T, output string)
	}{
		{
			name:   "replace git version",
			prompt: "This is a test: $(git --version)",
			validate: func(t *testing.T, output string) {
				assert.True(t,
					strings.Contains(output, "git version") || strings.Contains(output, "[git error:"),
					"expected git version or error, got '%s'", output,
				)
			},
		},
		{
			name:   "pretty format",
			prompt: `Commit info: $(git log -n 1 --pretty=format:"%h %s")`,
			validate: func(t *testing.T, output string) {
				assert.False(t,
					strings.Contains(output, "[error executing git command]"),
					"unexpected git error: '%s'", output,
				)
				assert.True(t,
					strings.HasPrefix(output, "Commit info: "),
					"expected output to start with 'Commit info: ', got '%s'", output,
				)
				assert.False(t,
					strings.Contains(output, "$(git") || strings.Contains(output, "%h"),
					"command substitution failed: got '%s'", output,
				)
				fmt.Println("Pretty output:", output)
			},
		},
		{
			name:   "git command error",
			prompt: "$(git not-a-real-command)",
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "[error executing git command]", "expected git error")
				fmt.Println(output)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := ExecuteGitCommands(tc.prompt)
			tc.validate(t, output)
		})
	}
}

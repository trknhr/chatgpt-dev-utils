# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ChatGPT Dev Utils is a CLI tool that sends prompts to ChatGPT via a Chrome extension. It consists of:
- A Go-based CLI application with an interactive TUI built using Bubble Tea
- A Chrome extension that injects prompts into the ChatGPT web interface
- WebSocket communication between CLI and extension

## Common Development Commands

### CLI Development (Go)
```bash
cd cli

# Run the CLI directly
go run .

# Run with extension support (starts WebSocket server)
go run . --extension

# Build the binary
make build  # or: go build -o cgpt

# Install globally
make install  # or: go install

# Debug with Delve
make debug

# Clean build artifacts
make clean
```

### Chrome Extension Development
1. Load extension in Chrome:
   - Navigate to `chrome://extensions/`
   - Enable Developer mode
   - Click "Load unpacked" and select the `extension/` directory
2. The extension automatically connects to the CLI's WebSocket server on `ws://localhost:8090`

## Architecture

### CLI Tool (`cli/`)
- **Entry point**: `cli/main.go`
- **UI Framework**: Bubble Tea (charmbracelet/bubbletea) for terminal UI
- **WebSocket**: Uses gorilla/websocket for extension communication
- **Key features**:
  - Interactive file selection and navigation
  - Built-in prompt templates (code review, commit messages, documentation)
  - Clipboard integration via atotto/clipboard

### Chrome Extension (`extension/`)
- **Manifest V3** extension
- **Components**:
  - `background.js`: Service worker managing WebSocket connection
  - `content.js`: Injects prompts into ChatGPT interface
  - `popup.html/js`: Extension popup UI
- **Permissions**: scripting, tabs, storage, host access to chatgpt.com

### Communication Flow
1. CLI starts WebSocket server on port 8090 when run with `--extension`
2. Extension connects to WebSocket server
3. CLI sends prompts through WebSocket
4. Extension injects prompts into ChatGPT's textarea and triggers submission

## Key Dependencies

### Go Dependencies
- `github.com/gorilla/websocket` - WebSocket server
- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/atotto/clipboard` - Clipboard operations

## Development Notes

- No test files currently exist in the project
- The CLI uses Go 1.24.2
- Chrome extension uses vanilla JavaScript (no framework)
- WebSocket server defaults to port 8090
- The project avoids OpenAI API costs by using the ChatGPT web interface directly
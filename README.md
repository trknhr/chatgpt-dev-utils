# 🧠 ChatGPT Dev Utils (`cdev`)

<p align="center">
  <img src="img/image.png" alt="ChatGPT Dev Utils Icon" width="160" />
</p>

A developer-friendly CLI tool to send prompts directly to ChatGPT using your Chrome session. Perfect for commit messages, code review, explanations, and more — without using the OpenAI API.

---

## ✨ Features

- Send prompts from CLI → Chrome extension → ChatGPT tab
- Built-in templates: code review, commit messages, documentation
- Interactive TUI interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- WebSocket-based integration (no OpenAI API key required)
- File/folder navigation & Git integration

---

## 📦 Installation (macOS / Linux)

Install the latest release via:

```bash
curl -sSfL https://raw.githubusercontent.com/your-org/chatgpt-dev-utils/main/install.sh | sh
```

This will:
- Detect your OS and CPU architecture
- Download the correct binary
- Install it to `/usr/local/bin/cdev`

---

## 🚀 Quick Start

### Interactive Mode

```bash
cdev
```

### With Chrome Extension

```bash
cdev --extension
```

You will be guided through:

1. Choosing prompt type (file / git)
2. Selecting files or Git templates
3. Editing prompt if needed
4. Copying prompt or sending to ChatGPT tab


## 🔌 Chrome Extension Setup

1. Open `chrome://extensions`
2. Load `extension/` directory as an unpacked extension
3. Open `chat.openai.com`
4. Ensure permissions allow WebSocket access


## 🧠 How It Works

```
┌────────────┐        ┌────────────────┐        ┌─────────────┐
│    cdev    ├──────▶│ Chrome Extension├──────▶│ chat.openai │
└────────────┘ WebSocket       │ Inject Prompt│
                               └─────────────┘
```

No OpenAI API keys. Works by controlling ChatGPT via browser.


## 🛠 Development

To run locally:

```bash
cd cli
go run . --extension
```

To build:

```bash
cd cli
go build -o cdev
./cdev --extension
```

Requires Go 1.24+

## 🧩 Templates Included

- Code Review (git diff)
- Commit Message
- Change Summary
- File Review
- Documentation

All templates are editable via TUI.

## 📬 Feedback & Contributions

PRs and issues welcome → [github.com/your-org/chatgpt-dev-utils](https://github.com/your-org/chatgpt-dev-utils)

## 📄 License

Apache 2.0 License — © 2025 Teruo Kunihiro


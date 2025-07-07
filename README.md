# 🧠 ChatGPT Dev Utils (`cdev`)

<p align="center">
  <img src="img/image.png" alt="ChatGPT Dev Utils Icon" width="160" />
</p>

A developer-friendly CLI tool to send prompts directly to ChatGPT using your Chrome session. Perfect for commit messages, code review, explanations, and more — without using the OpenAI API.

## ✨ Features

- Send prompts from CLI → Chrome extension → ChatGPT tab
- Built-in templates: code review, commit messages, documentation
- Interactive TUI interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- WebSocket-based integration (no OpenAI API key required)
- File/folder navigation & Git integration
- Official Chrome extension available on the [Chrome Web Store](https://chromewebstore.google.com/detail/chatgpt-dev-utils-extensi/bdfinimpohfncpgeokmamgfebfhnkebi)

## 📦 Installation (macOS / Linux)

Install the latest release via:

```bash
curl -sSfL https://raw.githubusercontent.com/trknhr/chatgpt-dev-utils/main/install.sh | sh
```

This will:
- Detect your OS and CPU architecture
- Download the correct binary
- Install it to `/usr/local/bin/cdev`

## 🚀 Quick Start

1. **Install the Chrome Extension** from the [Chrome Web Store](https://chromewebstore.google.com/detail/chatgpt-dev-utils-extensi/bdfinimpohfncpgeokmamgfebfhnkebi)
2. **Run the CLI tool:**

```bash
cdev
```

The WebSocket server starts automatically to enable Chrome extension integration.

You will be guided through:

1. Choosing prompt type (file / git)
2. Selecting files or Git templates
3. Editing prompt if needed
4. Copying prompt or sending to ChatGPT tab


## 🔌 Chrome Extension Setup

### Option 1: Install from Chrome Web Store (Recommended)

Install the official extension from the Chrome Web Store:

👉 **[ChatGPT Dev Utils Extension](https://chromewebstore.google.com/detail/chatgpt-dev-utils-extensi/bdfinimpohfncpgeokmamgfebfhnkebi)**

This is the easiest way to get started - just click "Add to Chrome" and you're ready to go!

### Option 2: Load Unpacked Extension (Development)

For development or if you prefer to load the extension manually:

1. Open `chrome://extensions`
2. Enable "Developer mode" 
3. Click "Load unpacked" and select the `extension/` directory
4. Open `chat.openai.com`
5. Ensure permissions allow WebSocket access

### Upgrading from Unpacked to Chrome Web Store Version

If you're currently using the unpacked extension:

1. Remove the unpacked extension from `chrome://extensions`
2. Install from the [Chrome Web Store](https://chromewebstore.google.com/detail/chatgpt-dev-utils-extensi/bdfinimpohfncpgeokmamgfebfhnkebi)
3. The Chrome Web Store version will automatically update with new features and bug fixes


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
go run .
```

To build:

```bash
cd cli
go build -o cdev
./cdev
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

PRs and issues welcome → [github.com/trknhr/chatgpt-dev-utils](https://github.com/trknhr/chatgpt-dev-utils)

## 📄 License

Apache 2.0 License — © 2025 Teruo Kunihiro

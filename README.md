# ChatGPT Dev Utils

<p align="center">
  <img src="img/image.png" alt="ChatGPT Dev Utils Icon" width="160" />
</p>

A developer-friendly CLI tool that lets you send prompts to ChatGPT directly from your terminal, using your existing ChatGPT web session via a Chrome extension. Useful for generating commit messages, code explanations, translations, and more.

## ✨ Features

- Send prompts from CLI → Chrome extension → ChatGPT tab
- Built-in prompt templates for code review, commit messages, and documentation
- Interactive CLI mode for selecting files and templates
- Works via WebSocket server (Go-based)
- Avoids using OpenAI API (cost-effective)

## 📦 Installation

You need Go 1.21+ installed.

```bash
cd cli
# Run interactively (no build):
go run .
# Or build and use the binary:
go build -o cgpt
./cgpt
```

## 🧪 Usage

### Interactive mode:
```bash
cd cli
go run .
# or if built:
./cgpt
```

- Choose prompt type: File-based or Git-based
- Select files or git template
- Copy prompt to clipboard or send to extension (if enabled)

### With Chrome Extension:
```bash
cd cli
go run . --extension
# or
./cgpt --extension
```

## 🧩 Templates

The CLI provides built-in templates for common tasks:
- **File-based**: Code Review, Documentation, Custom
- **Git-based**: Code Review, Commit Message, Change Summary, Custom

You can select and edit these templates interactively. (User-defined presets are not currently supported.)

---

## 🧠 How It Works

1. CLI (Go) starts a local WebSocket server (`--extension` flag)
2. Chrome extension connects to it via WebSocket
3. CLI sends prompt to extension
4. Extension injects the prompt into ChatGPT tab
5. You get a response, visible in browser

---

## 🛠 Development

```bash
cd cli
go run .
# or build:
go build -o cgpt
./cgpt
```

---

## 🧩 Extension Setup

1. Load `extension/` into Chrome as an unpacked extension
2. Ensure `chatgpt.com` is open
3. Allow extension permissions

---

## 🤝 License

This project is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).
Copyright © 2025 Teruo Kunihiro.

# 🔍 Glimpse - Interactive Code Search

A lightning-fast, beautiful terminal-based code search tool with real-time results and instant editor integration.

<img width="971" height="374" alt="Screenshot 2025-07-24 at 9 07 17 PM" src="https://github.com/user-attachments/assets/067860aa-4c12-45b7-8609-79ea3efc8b40" />

## ✨ Features

- **🚀 Real-time Search** - Results update as you type
- **🎨 Beautiful TUI** - Modern, colorful interface with split-pane layout
- **⚡ Lightning Fast** - Optimized search with smart file filtering
- **🎯 Instant Editor Integration** - Press Enter to open files in VS Code, Vim, or your preferred editor
- **📱 Responsive Design** - Clean, compact interface that works in any terminal size
- **🔧 Smart Filtering** - Automatically skips binary files and build directories
- **⌨️ Vim-like Navigation** - Use ctrl+j/ctrl+k or arrow keys for navigation
- **🔤 Case Toggle** - Ctrl+I to switch between case-sensitive and case-insensitive search
- **🧹 Auto-cleanup** - Clears terminal after use for a clean workspace
  
  heavily inspired by https://github.com/nvim-telescope/telescope.nvim

## 🎯 Why Glimpse?

Traditional tools like `grep` and `find` are powerful but:
- ❌ Results are static and hard to navigate
- ❌ No preview of file contents
- ❌ Require complex command syntax
- ❌ Don't integrate with modern editors

**Glimpse solves this with:**
- ✅ Interactive, real-time search
- ✅ File preview with syntax context
- ✅ One-key editor opening
- ✅ Beautiful, intuitive interface

## 🚀 Installation

### Prerequisites
- Go 1.21 or later

### Install from source
```bash
git clone https://github.com/pixelknightdev/Glimpse.git
cd glimpse
go build -o glimpse cmd/main.go
cp glimpse $(go env GOPATH)/bin/
```

### Install directly
```bash
go install github.com/pixelknightdev/Glimpse/cmd/gimpse@latest
```

## 📖 Usage

### Interactive Mode (Default)
```bash
# Launch interactive search
glimpse

# Then:
# - Type to search in real-time
# - Use ↑/↓ or ctrl+j/ctrl+k to navigate results
# - Press Enter to open file in editor
# - Press ctrl+c to quit
```

### CLI Mode
```bash
# Traditional grep-like output
glimpse --cli "search term"

# Case-insensitive search
glimpse --cli -i "search term"
```

## ⌨️ Keyboard Controls

| Key | Action |
|-----|--------|
| `Type` | Search in real-time |
| `↑/↓` or `ctrl+j/ctrl+k` | Navigate through results |
| `Enter` | Open selected file in editor |
| `Ctrl+I` | Toggle case sensitivity |
| `Ctrl+C` | Quit |

## 🎨 Video demo



https://github.com/user-attachments/assets/48c9256d-5faf-45b8-b541-4ab7fbe473a3



## 🛠 Technical Details

### Architecture
```
glimpse/
├── cmd/main.go              # CLI interface and mode routing
├── internal/
│   ├── search/             # Core search engine
│   │   └── search.go      # File traversal and pattern matching
│   └── tui/               # Terminal user interface
│       ├── model.go       # Bubbletea TUI implementation
│       └── editor.go      # Cross-platform editor integration
├── go.mod                 # Dependencies
└── README.md
```

### Performance Optimizations
- **Concurrent file processing** with goroutines
- **Smart binary file detection** by extension and content
- **Result limiting** to prevent memory issues
- **Directory exclusion** for common build/cache folders
- **Early termination** when sufficient results found

### Supported Editors
- **VS Code** (`code` command)
- **Vim/Neovim** 
- **System default editor**
- Cross-platform support (macOS, Linux, Windows)

## 🎯 Use Cases

- **🔍 Code exploration** - Quickly find function definitions, imports, or patterns
- **🐛 Debugging** - Locate error messages, variable usages, or specific logic
- **📚 Learning codebases** - Navigate unfamiliar projects with ease
- **🔧 Refactoring** - Find all instances of code that needs updating
- **📖 Documentation** - Search for comments, TODOs, or documentation strings

## 🤝 Contributing

We welcome contributions! Here are some ways to help:

1. **🐛 Bug Reports** - Found an issue? Open a GitHub issue
2. **💡 Feature Requests** - Have an idea? Let's discuss it
3. **🛠 Code Contributions** - Submit a pull request
4. **📖 Documentation** - Help improve our docs
5. **🌟 Spread the word** - Star the repo and share with others

### Development Setup
```bash
git clone https://github.com/pixelknightdev/Glimpse.git
cd glimpse
go mod tidy
go run cmd/main.go
```

## 📋 Roadmap

- [ ] **📦 Package managers** - Homebrew, apt, chocolatey
- [ ] **🎨 Themes** - Customizable color schemes
- [ ] **🔧 Config files** - User preferences and settings
- [ ] **📊 Regex support** - Advanced pattern matching
- [ ] **📁 File filtering** - Include/exclude by file type
- [ ] **🔗 Git integration** - Search only modified files
- [ ] **💾 Search history** - Remember recent searches
- [ ] **🚀 Fuzzy search** - More flexible matching

## 📊 Comparison

| Tool | Real-time | Interactive | Editor Integration | Preview | Performance |
|------|-----------|-------------|-------------------|---------|-------------|
| **Glimpse** | ✅ | ✅ | ✅ | ✅ | ⚡ |
| grep | ❌ | ❌ | ❌ | ❌ | ⚡ |
| ripgrep | ❌ | ❌ | ❌ | ❌ | ⚡⚡ |
| fzf | ❌ | ✅ | ⚠️ | ❌ | ⚡ |
| VS Code Search | ❌ | ✅ | ✅ | ✅ | 🐌 |

## 🙏 Acknowledgments

Built with these amazing tools:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=pixelknightdev/glimpse&type=Date)](https://star-history.com/#pixelknightdev/glimpse&Date)

---

<div align="center">

**[⭐ Star us on GitHub](https://github.com/pixelknightdev/glimpse)** • **[🐛 Report Issues](https://github.com/pixelknightdev/glimpse/issues)** • **[💡 Request Features](https://github.com/pixelknightdev/glimpse/issues)**

Made with ❤️ by developers, for developers

</div>

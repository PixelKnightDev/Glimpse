package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixelknightdev/glimpse/internal/search"
	"github.com/pixelknightdev/glimpse/internal/tui"
)

func main() {
	var caseInsensitive = flag.Bool("i", false, "Case insensitive search")
	var help = flag.Bool("h", false, "Show help")
	var cliMode = flag.Bool("cli", false, "Force CLI mode instead of interactive TUI")

	flag.Parse()

	if *help {
		fmt.Println("🔍 Glimpse - Interactive Code Search")
		fmt.Println("\nUsage:")
		fmt.Println("  glimpse                    # Interactive TUI mode (default)")
		fmt.Println("  glimpse --cli <term>       # CLI mode")
		fmt.Println("  glimpse -i <term>          # Case insensitive CLI search")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nTUI Controls:")
		fmt.Println("  Type to search, ↑/↓ navigate, Enter to open, q to quit")
		return
	}

	args := flag.Args()

	// If CLI mode is explicitly requested AND has search term
	if *cliMode && len(args) > 0 {
		searchTerm := args[0]
		options := search.SearchOptions{
			CaseInsensitive: *caseInsensitive,
			MaxResults:      0, // No limit for CLI mode
		}

		results := search.SearchFiles(searchTerm, ".", options)

		for _, result := range results {
			fmt.Printf("%s:%d: %s\n", result.File, result.Line, result.Content)
		}

		if len(results) == 0 {
			fmt.Println("No matches found")
		}
		return
	}

	// Default: Launch interactive TUI mode
	model := tui.InitialModel()
	p := tea.NewProgram(model)
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	
	// Clear the terminal after TUI exits
	clearTerminal()
}

// Clear terminal screen across different operating systems
func clearTerminal() {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "linux", "darwin": // Linux and macOS
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		// Fallback: print ANSI escape sequence
		fmt.Print("\033[2J\033[H")
		return
	}
	
	cmd.Stdout = os.Stdout
	cmd.Run()
}

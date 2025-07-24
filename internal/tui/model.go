package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pixelknightdev/glimpse/internal/search"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#F25D94"))

	leftPanelStyle = lipgloss.NewStyle().
			Width(50).
			Height(10).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FAFAFA")).
			Padding(1)

	rightPanelStyle = lipgloss.NewStyle().
			Width(60).
			Height(10).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FAFAFA")).
			Padding(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Align(lipgloss.Center).
			Width(115)

	contextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	matchLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F25D94")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
)

const visibleItems = 6

type Model struct {
	searchInput     textinput.Model
	results         []search.Result
	selectedIndex   int
	scrollOffset    int
	preview         viewport.Model
	caseInsensitive bool
	lastMessage     string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type to search your code..."
	ti.Focus()
	ti.Width = 50

	vp := viewport.New(58, 8)
	vp.SetContent("🚀 Start typing to search!\n\n✨ Real-time results will appear here")

	return Model{
		searchInput:     ti,
		results:         []search.Result{},
		preview:         vp,
		caseInsensitive: true,
		scrollOffset:    0,
		lastMessage:     "",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+i":
			m.caseInsensitive = !m.caseInsensitive
			statusText := "CASE-INSENSITIVE"
			if !m.caseInsensitive {
				statusText = "CASE-SENSITIVE"
			}
			m.lastMessage = "🔧 Switched to " + statusText + " mode"
			if query := m.searchInput.Value(); query != "" {
				options := search.SearchOptions{CaseInsensitive: m.caseInsensitive, MaxResults: 50}
				m.results = search.SearchFiles(query, ".", options)
				m.selectedIndex = 0
				m.scrollOffset = 0
				m.updatePreview()
			}
		case "up", "ctrl+k":
			if len(m.results) > 0 && m.selectedIndex > 0 {
				m.selectedIndex--
				m.adjustScroll()
				m.updatePreview()
				m.lastMessage = ""
			}
		case "down", "ctrl+j":
			if len(m.results) > 0 && m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
				m.adjustScroll()
				m.updatePreview()
				m.lastMessage = ""
			}
		case "enter":
			if len(m.results) > 0 {
				selected := m.results[m.selectedIndex]
				OpenFileInEditor(selected.File, selected.Line)
				return m, tea.Quit
			}
		default:
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)

			query := m.searchInput.Value()
			if query != "" {
				options := search.SearchOptions{CaseInsensitive: m.caseInsensitive, MaxResults: 50}
				m.results = search.SearchFiles(query, ".", options)
				m.selectedIndex = 0
				m.scrollOffset = 0
				m.updatePreview()
				m.lastMessage = ""
			} else {
				m.results = []search.Result{}
				m.scrollOffset = 0
				m.preview.SetContent("🚀 Start typing to search!\n\n✨ Real-time results will appear here")
				m.lastMessage = ""
			}

			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) adjustScroll() {
	if m.selectedIndex >= m.scrollOffset+visibleItems {
		m.scrollOffset = m.selectedIndex - visibleItems + 1
	}

	if m.selectedIndex < m.scrollOffset {
		m.scrollOffset = m.selectedIndex
	}

	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m *Model) updatePreview() {
	if len(m.results) == 0 {
		m.preview.SetContent("❌ No matches found\n\nTry a different search term!")
		return
	}

	selected := m.results[m.selectedIndex]
	contextLines := m.getFileContext(selected.File, selected.Line, 1)

	content := "📁 " + selected.File + "\n"
	content += "📍 Line " + fmt.Sprintf("%d", selected.Line) + "\n\n"

	if len(contextLines) > 0 {
		content += "🔍 GLIMPSE:\n"
		content += strings.Repeat("─", 25) + "\n"

		startLine := selected.Line - 1
		if startLine < 1 {
			startLine = 1
		}

		for i, line := range contextLines {
			lineNum := startLine + i
			if lineNum == selected.Line {
				content += matchLineStyle.Render(fmt.Sprintf("%3d | %s", lineNum, line)) + "\n"
			} else {
				content += contextStyle.Render(fmt.Sprintf("%3d | %s", lineNum, line)) + "\n"
			}
		}
		content += strings.Repeat("─", 25) + "\n"
	} else {
		content += "🔍 MATCH:\n"
		content += strings.Repeat("─", 25) + "\n"
		content += selected.Content + "\n"
		content += strings.Repeat("─", 25) + "\n"
	}

	content += "\n💡 Press ENTER to open & close"
	content += "\n" + fmt.Sprintf("Result %d of %d", m.selectedIndex+1, len(m.results))

	m.preview.SetContent(content)
}

func (m *Model) getFileContext(filename string, targetLine int, contextSize int) []string {
	file, err := os.Open(filename)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNum := 1

	startLine := targetLine - contextSize
	if startLine < 1 {
		startLine = 1
	}
	endLine := targetLine + contextSize

	for scanner.Scan() {
		if lineNum >= startLine && lineNum <= endLine {
			lines = append(lines, scanner.Text())
		}
		if lineNum > endLine {
			break
		}
		lineNum++
	}

	return lines
}

func (m Model) View() string {
	caseStatus := "CASE-INSENSITIVE"
	if !m.caseInsensitive {
		caseStatus = "CASE-SENSITIVE"
	}

	title := titleStyle.Render("🔍 GLIMPSE - Interactive Code Search")
	status := statusStyle.Render("Mode: " + caseStatus)

	messageText := ""
	if m.lastMessage != "" {
		messageText = "  " + successStyle.Render(m.lastMessage)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", status, messageText)

	searchBox := "🔎 Glimpse: " + m.searchInput.View()

	leftContent := m.renderResults()
	leftPanel := leftPanelStyle.Render(leftContent)

	rightPanel := rightPanelStyle.Render(m.preview.View())

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	)

	help := helpStyle.Render("(↑/↓)/(ctrl+k/ctrl+j): Navigate   Enter: Open & Close   ctrl+c: Quit   Ctrl+I: Toggle Case")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		searchBox,
		"",
		mainContent,
		"",
		help,
	)
}

func (m Model) renderResults() string {
	if len(m.results) == 0 {
		return "💭 No results yet...\n\nStart typing to search!"
	}

	lines := []string{}

	totalResults := len(m.results)
	startIdx := m.scrollOffset + 1
	endIdx := m.scrollOffset + visibleItems
	if endIdx > totalResults {
		endIdx = totalResults
	}

	lines = append(lines, fmt.Sprintf("📊 Results %d-%d of %d:", startIdx, endIdx, totalResults))
	lines = append(lines, "")

	for i := 0; i < visibleItems && (m.scrollOffset+i) < len(m.results); i++ {
		resultIdx := m.scrollOffset + i
		result := m.results[resultIdx]
		line := fmt.Sprintf("%s:%d", result.File, result.Line)

		if len(line) > 45 {
			line = line[:42] + "..."
		}

		if resultIdx == m.selectedIndex {
			line = selectedStyle.Render("▶ " + line)
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	if m.scrollOffset > 0 {
		lines[1] = statusStyle.Render("  ↑ More above ↑")
	}

	if m.scrollOffset+visibleItems < len(m.results) {
		lines = append(lines, statusStyle.Render("  ↓ More below ↓"))
	}

	return strings.Join(lines, "\n")
}

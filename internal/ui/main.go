package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ketsu/ketsu/internal/core"
	"github.com/ketsu/ketsu/internal/editor"
	"github.com/ketsu/ketsu/internal/finder"
	"github.com/ketsu/ketsu/internal/markdown"
)

type model struct {
	core           *core.Core
	theme          *ThemeManager
	markdown       *markdown.Renderer
	finder         *finder.Finder
	width          int
	height         int
	commandInput   string
	shellInput     string
	shellMode      bool
	searchMode     bool
	searchInput    string
	showExplorer   bool
	showPreview    bool
	finderMode     bool    // Fuzzy finder mode
	finderType     string  // "file", "tag", "search"
	
	selectedFile   int
	explorerHeight int
	editorScroll   int
	previewScroll  int
	
	headerHeight   int
	footerHeight   int
	
	// File path display
	filePath       string
	
	// Pending command for visual mode
	pendingCmd     string
}

func New() *model {
	return &model{
		core:           core.New(),
		theme:          NewThemeManager(),
		markdown:       markdown.New(),
		finder:         finder.New(),
		showExplorer:   true,
		showPreview:    false,
		finderMode:     false,
		finderType:     "file",
		selectedFile:   0,
		headerHeight:   3,
		footerHeight:   3,
	}
}

func (m *model) Init() tea.Cmd {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	
	defaultSpace := filepath.Join(home, "notes")
	
	if _, err := os.Stat(".ketsu"); err == nil {
		defaultSpace, _ = os.Getwd()
	}
	
	if err := m.core.LoadSpace(defaultSpace); err != nil {
		m.core.StatusMessage = fmt.Sprintf("Error loading space: %v", err)
	}
	
	// Update file path in UI
	if m.core.Editor.FilePath != "" {
		m.filePath = m.core.Editor.FilePath
	}
	
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	
	// Update file path if changed
	if m.core.Editor.FilePath != m.filePath {
		m.filePath = m.core.Editor.FilePath
	}
	
	return m, nil
}

func (m *model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.finderMode {
		return m.handleFinderKey(msg)
	}
	
	if m.shellMode {
		return m.handleShellKey(msg)
	}
	
	if m.searchMode {
		return m.handleSearchKey(msg)
	}
	
	if m.core.Editor.Mode == editor.ModeCommand {
		return m.handleCommandKey(msg)
	}
	
	if m.core.Editor.Mode == editor.ModeSearch {
		return m.handleEditorSearchKey(msg)
	}
	
	if m.core.Editor.Mode == editor.ModeVisual || 
	   m.core.Editor.Mode == editor.ModeVisualLine ||
	   m.core.Editor.Mode == editor.ModeVisualBlock {
		return m.handleVisualKey(msg)
	}
	
	if m.core.CurrentView == core.ViewEditor {
		return m.handleEditorKey(msg)
	}
	
	return m.handleExplorerKey(msg)
}

func (m *model) handleKeyCommon(msg tea.KeyMsg) {
	switch msg.String() {
	case "ctrl+c", "ctrl+q":
	case ":":
		m.core.Editor.EnterCommandMode()
	case "/":
		m.core.Editor.EnterSearchMode()
	case "?":
		m.searchMode = true
		m.searchInput = ""
	case "\\":
		m.shellMode = true
		m.shellInput = ""
	case "ctrl+s":
		m.core.SaveKV()
		m.core.SaveFile()
		m.core.StatusMessage = "Saved"
	case "ctrl+f":
		// Open file finder
		m.openFinder("file")
	case "ctrl+g":
		// Open search finder
		m.openFinder("search")
	case "tab":
		m.showExplorer = !m.showExplorer
	case "shift+tab":
		m.showPreview = !m.showPreview
	case "ctrl+h":
		m.showExplorer = !m.showExplorer
	case "ctrl+p":
		m.showPreview = !m.showPreview
	case "ctrl+shift+t":
		// Cycle through themes
		themes := m.theme.GetAvailableThemes()
		current := m.theme.GetName()
		nextIdx := 0
		for i, t := range themes {
			if t == current {
				nextIdx = (i + 1) % len(themes)
				break
			}
		}
		m.theme.SetTheme(themes[nextIdx])
		m.core.StatusMessage = fmt.Sprintf("Theme: %s", themes[nextIdx])
	}
}

func (m *model) handleExplorerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.handleKeyCommon(msg)
	
	dir := m.core.SpaceDir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	
	entries, _ := os.ReadDir(dir)
	var visibleEntries []os.DirEntry
	for _, e := range entries {
		if len(e.Name()) > 0 && e.Name()[0] != '.' {
			visibleEntries = append(visibleEntries, e)
		}
	}
	
	maxSelected := len(visibleEntries) - 1
	if maxSelected < 0 {
		maxSelected = 0
	}
	
	switch msg.String() {
	case "j", "down", "ctrl+n":
		if m.selectedFile < maxSelected {
			m.selectedFile++
		}
	case "k", "up":
		if m.selectedFile > 0 {
			m.selectedFile--
		}
	case "enter":
		if m.selectedFile < len(visibleEntries) {
			filename := visibleEntries[m.selectedFile].Name()
			path := filepath.Join(dir, filename)
			
			if visibleEntries[m.selectedFile].IsDir() {
				// Enter directory
				m.core.SpaceDir = path
				m.selectedFile = 0
			} else {
				m.core.OpenFile(path)
				// Auto show preview for markdown files
				if strings.HasSuffix(strings.ToLower(filename), ".md") {
					m.showPreview = true
				}
			}
		}
	case "l", "right":
		m.core.CurrentView = core.ViewEditor
	case "h", "left":
		// Go up one directory
		parent := filepath.Dir(m.core.SpaceDir)
		if parent != m.core.SpaceDir {
			m.core.SpaceDir = parent
			m.selectedFile = 0
		}
	case "g":
		m.selectedFile = 0
	case "G":
		m.selectedFile = maxSelected
	case "n":
		// New file
		m.core.Editor.Load("")
		m.core.Editor.FilePath = ""
		m.core.CurrentView = core.ViewEditor
		m.core.StatusMessage = "New file"
	}
	
	return m, nil
}

func (m *model) handleEditorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.handleKeyCommon(msg)
	
	e := m.core.Editor
	
	// Undo/Redo
	if msg.String() == "u" {
		if e.Undo() {
			m.core.StatusMessage = "Undo"
		}
		return m, nil
	}
	if msg.String() == "ctrl+r" {
		if e.Redo() {
			m.core.StatusMessage = "Redo"
		}
		return m, nil
	}
	
	// Visual mode
	if msg.String() == "v" {
		e.EnterVisualMode()
		m.pendingCmd = ""
		return m, nil
	}
	if msg.String() == "V" {
		e.EnterVisualLineMode()
		m.pendingCmd = ""
		return m, nil
	}
	
	switch msg.String() {
	case "esc":
		if e.Mode == editor.ModeInsert {
			e.Mode = editor.ModeNormal
		}
	case "i":
		e.Mode = editor.ModeInsert
	case "a":
		e.MoveRight()
		e.Mode = editor.ModeInsert
	case "A":
		e.MoveToLineEnd()
		e.Mode = editor.ModeInsert
	case "I":
		e.MoveToFirstNonBlank()
		e.Mode = editor.ModeInsert
	case "R":
		e.EnterReplaceMode()
	case "o":
		e.OpenLine()
		e.Mode = editor.ModeInsert
	case "O":
		e.OpenLineAbove()
		e.Mode = editor.ModeInsert
	case "x", "delete":
		e.DeleteCharForward()
	case "X":
		e.DeleteChar()
	case "backspace":
		e.MoveLeft()
		e.DeleteChar()
	case "0":
		e.MoveToLineStart()
	case "^":
		e.MoveToFirstNonBlank()
	case "$":
		e.MoveToLineEnd()
	case "w":
		e.MoveWordForward()
	case "W":
		for i := 0; i < 5; i++ {
			e.MoveWordForward()
		}
	case "b":
		e.MoveWordBackward()
	case "e":
		e.MoveToEndOfWord()
	case "h", "left":
		e.MoveLeft()
	case "l", "right":
		e.MoveRight()
	case "j", "down":
		e.MoveDown()
	case "k", "up":
		e.MoveUp()
	case "gg":
		e.CursorY = 0
		e.CursorX = 0
	case "G":
		e.CursorY = len(e.Lines) - 1
		if e.CursorX > len(e.Lines[e.CursorY]) {
			e.CursorX = len(e.Lines[e.CursorY])
		}
	case "ctrl+g":
		m.core.CurrentView = core.ViewEditor
		m.showExplorer = true
	case "~":
		e.ToggleCase()
	case "J":
		e.JoinLines()
	case ">":
		e.IndentLine()
	case "<":
		e.OutdentLine()
	}
	
	// Multi-key commands (dd, yy, cc, etc.)
	if msg.String() == "d" {
		if m.pendingCmd == "d" {
			e.DeleteLine()
			m.core.StatusMessage = "Line deleted"
			m.pendingCmd = ""
		} else {
			m.pendingCmd = "d"
		}
		return m, nil
	}
	if msg.String() == "c" {
		if m.pendingCmd == "c" {
			e.ChangeLine()
			m.core.StatusMessage = "Line changed"
			m.pendingCmd = ""
		} else {
			m.pendingCmd = "c"
		}
		return m, nil
	}
	if msg.String() == "y" {
		if m.pendingCmd == "y" {
			e.YankLine()
			m.core.StatusMessage = "Line yanked"
			m.pendingCmd = ""
		} else {
			m.pendingCmd = "y"
		}
		return m, nil
	}
	
	// Clear pending command if another key is pressed
	if m.pendingCmd != "" && msg.String() != "d" && msg.String() != "c" && msg.String() != "y" {
		m.pendingCmd = ""
	}
	
	// Search commands
	if msg.String() == "n" {
		if e.SearchNext() {
			m.core.StatusMessage = "Next match"
		} else {
			m.core.StatusMessage = "No more matches"
		}
	}
	if msg.String() == "N" {
		if e.SearchPrevious() {
			m.core.StatusMessage = "Previous match"
		} else {
			m.core.StatusMessage = "No more matches"
		}
	}
	if msg.String() == "*" {
		// Search word under cursor
		line := e.GetCurrentLine()
		word := extractWordAt(line, e.CursorX)
		if word != "" {
			e.Search(word)
			m.core.StatusMessage = "Searching: " + word
		}
	}
	
	// Paste commands
	if msg.String() == "p" {
		e.PasteAfter()
	}
	if msg.String() == "P" {
		e.PasteBefore()
	}
	
	// Replace
	if msg.String() == "r" {
		// Next character will be the replacement (single char)
		m.pendingCmd = "r"
		return m, nil
	}
	
	// Handle pending commands
	if m.pendingCmd == "r" {
		for _, ch := range msg.Runes {
			e.ReplaceChar(ch)
			m.pendingCmd = ""
			break
		}
		return m, nil
	}
	
	// Handle Insert mode
	if e.Mode == editor.ModeInsert {
		for _, ch := range msg.Runes {
			e.InsertChar(ch)
		}
		return m, nil
	}
	
	// Handle Replace mode
	if e.Mode == editor.ModeReplace {
		for _, ch := range msg.Runes {
			e.ReplaceChar(ch)
		}
		return m, nil
	}
	
	return m, nil
}

func (m *model) handleVisualKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	e := m.core.Editor
	
	switch msg.String() {
	case "esc", "v", "V":
		e.ExitVisualMode()
		return m, nil
	case "d", "x":
		e.DeleteSelection()
		m.core.StatusMessage = "Selection deleted"
		return m, nil
	case "y", "Y":
		e.YankSelection()
		m.core.StatusMessage = "Selection yanked"
		e.ExitVisualMode()
		return m, nil
	case "c":
		e.YankSelection()
		e.DeleteSelection()
		e.Mode = editor.ModeInsert
		m.core.StatusMessage = "Selection changed"
		return m, nil
	case "h", "left":
		e.MoveLeft()
		e.UpdateVisualSelection()
	case "l", "right":
		e.MoveRight()
		e.UpdateVisualSelection()
	case "j", "down":
		e.MoveDown()
		e.UpdateVisualSelection()
	case "k", "up":
		e.MoveUp()
		e.UpdateVisualSelection()
	case "w":
		e.MoveWordForward()
		e.UpdateVisualSelection()
	case "b":
		e.MoveWordBackward()
		e.UpdateVisualSelection()
	case "0":
		e.MoveToLineStart()
		e.UpdateVisualSelection()
	case "$":
		e.MoveToLineEnd()
		e.UpdateVisualSelection()
	case "gg":
		e.CursorY = 0
		e.CursorX = 0
		e.UpdateVisualSelection()
	case "G":
		e.CursorY = len(e.Lines) - 1
		e.UpdateVisualSelection()
	}
	
	return m, nil
}

func (m *model) handleEditorSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	e := m.core.Editor
	
	switch msg.String() {
	case "enter":
		pattern := e.CommandBuf
		e.ExitSearchMode()
		if pattern != "" {
			if e.Search(pattern) {
				m.core.StatusMessage = "Found: " + pattern
			} else {
				m.core.StatusMessage = "Not found: " + pattern
			}
		}
	case "backspace":
		if len(e.CommandBuf) > 0 {
			e.CommandBuf = e.CommandBuf[:len(e.CommandBuf)-1]
		}
	case "esc":
		e.ExitSearchMode()
	default:
		for _, ch := range msg.Runes {
			e.TypeSearch(ch)
		}
	}
	
	return m, nil
}

func (m *model) handleCommandKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	e := m.core.Editor
	
	switch msg.String() {
	case "enter":
		cmdLine := e.CommandLine
		e.CommandLine = ""
		e.Mode = editor.ModeNormal
		
		// Handle theme command
		if strings.HasPrefix(strings.ToLower(cmdLine), "theme") {
			parts := strings.Fields(cmdLine)
			if len(parts) > 1 {
				themeName := strings.ToLower(parts[1])
				if m.theme.SetTheme(themeName) {
					m.core.StatusMessage = fmt.Sprintf("Theme set to: %s", themeName)
				} else {
					m.core.StatusMessage = fmt.Sprintf("Unknown theme: %s. Available: %s", 
						themeName, strings.Join(m.theme.GetAvailableThemes(), ", "))
				}
			} else {
				m.core.StatusMessage = fmt.Sprintf("Current theme: %s. Available: %s", 
					m.theme.GetName(), strings.Join(m.theme.GetAvailableThemes(), ", "))
			}
			return m, nil
		}
		
		// Handle preview toggle
		if strings.ToLower(cmdLine) == "preview" || strings.ToLower(cmdLine) == "p" {
			m.showPreview = !m.showPreview
			m.core.StatusMessage = fmt.Sprintf("Preview: %v", m.showPreview)
			return m, nil
		}
		
		// Handle substitution command (:s/old/new/g)
		if strings.HasPrefix(cmdLine, "s/") || strings.HasPrefix(cmdLine, "substitute ") {
			count := m.handleSubstitute(cmdLine)
			if count > 0 {
				m.core.StatusMessage = fmt.Sprintf("%d substitution(s)", count)
			} else {
				m.core.StatusMessage = "Pattern not found"
			}
			return m, nil
		}
		
		// Handle other commands
		result := e.ExecuteCommand(cmdLine)
		
		switch result {
		case "quit_save":
			m.core.SaveKV()
			m.core.SaveFile()
			return m, tea.Quit
		case "quit", "quit_force":
			return m, tea.Quit
		case "write":
			m.core.SaveFile()
			m.core.StatusMessage = "File saved"
		case "write_force":
			m.core.SaveFile()
			m.core.SaveKV()
			m.core.StatusMessage = "All saved"
		}
		
		if len(result) > 3 && result[:3] == "kv:" {
			cmd := result[3:]
			result := m.core.ProcessCommand(cmd)
			m.core.StatusMessage = result
		}
		
	case "backspace":
		if len(e.CommandLine) > 0 {
			e.CommandLine = e.CommandLine[:len(e.CommandLine)-1]
		}
	case "esc":
		e.ExitCommandMode()
	default:
		for _, ch := range msg.Runes {
			e.TypeCommand(ch)
		}
	}
	
	return m, nil
}

func (m *model) handleShellKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		result := m.core.ProcessCommand(m.shellInput)
		if result != "" {
			m.core.StatusMessage = result
		}
		m.shellInput = ""
	case "backspace":
		if len(m.shellInput) > 0 {
			m.shellInput = m.shellInput[:len(m.shellInput)-1]
		}
	case "esc":
		m.shellMode = false
		m.shellInput = ""
	case "ctrl+c":
		m.shellMode = false
		m.shellInput = ""
	default:
		for _, ch := range msg.Runes {
			m.shellInput += string(ch)
		}
	}
	
	return m, nil
}

func (m *model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.searchMode = false
		m.core.SearchQuery = m.searchInput
	case "backspace":
		if len(m.searchInput) > 0 {
			m.searchInput = m.searchInput[:len(m.searchInput)-1]
		}
	case "esc":
		m.searchMode = false
		m.searchInput = ""
	case "ctrl+c":
		m.searchMode = false
		m.searchInput = ""
	default:
		for _, ch := range msg.Runes {
			m.searchInput += string(ch)
		}
	}
	
	return m, nil
}

func (m *model) View() string {
	t := m.theme.GetTheme()
	var s string
	
	// If finder mode is active, show finder overlay
	if m.finderMode {
		return m.viewFinder()
	}
	
	s += m.viewHeader()
	
	if m.core.CurrentView == core.ViewEditor {
		if m.showExplorer && m.showPreview {
			// Three column layout
			s += m.viewEditorWithExplorerAndPreview()
		} else if m.showExplorer {
			s += m.viewEditorWithExplorer()
		} else if m.showPreview {
			s += m.viewEditorWithPreview()
		} else {
			s += m.viewEditor()
		}
	} else {
		s += m.viewExplorer()
	}
	
	// Command line input
	if m.core.Editor.Mode == editor.ModeCommand {
		s += "\n" + t.StatusBarStyle().Render(":"+m.core.Editor.CommandLine+"_")
	}
	
	// Search input (editor search)
	if m.core.Editor.Mode == editor.ModeSearch {
		s += "\n" + t.StatusBarStyle().Render("/"+m.core.Editor.CommandBuf+"_")
	}
	
	// Shell input
	if m.shellMode {
		s += "\n" + t.StatusBarStyle().Render("> "+m.shellInput+"_")
	}
	
	// Global search input
	if m.searchMode {
		s += "\n" + t.StatusBarStyle().Render("/ "+m.searchInput+"_ (ESC to cancel)")
	}
	
	s += m.viewFooter()
	
	return s
}

func (m *model) viewHeader() string {
	t := m.theme.GetTheme()
	
	// File name and path
	fileName := "No file"
	if m.core.Editor.FilePath != "" {
		fileName = filepath.Base(m.core.Editor.FilePath)
	}
	
	// Mode indicator
	modeStr := m.core.Editor.GetModeString()
	
	// Header with file info
	header := t.HeaderStyle().Render("Ketsu") + " " +
		t.ButtonStyle().Render(fileName) + " " +
		t.ModeStyle(modeStr).Render(modeStr) + "\n"
	
	// Toolbar buttons
	expBtn := t.ButtonActiveStyle().Render("Files")
	if !m.showExplorer {
		expBtn = t.ButtonStyle().Render("Files")
	}
	
	previewBtn := t.ButtonActiveStyle().Render("Preview")
	if !m.showPreview {
		previewBtn = t.ButtonStyle().Render("Preview")
	}
	
	shellBtn := t.ButtonActiveStyle().Render("Shell")
	if !m.shellMode {
		shellBtn = t.ButtonStyle().Render("Shell")
	}
	
	searchBtn := t.ButtonActiveStyle().Render("Search")
	if !m.searchMode {
		searchBtn = t.ButtonStyle().Render("Search")
	}
	
	saveBtn := t.ButtonStyle().Render("Save")
	themeBtn := t.ButtonStyle().Render("Theme")
	
	toolbar := expBtn + " " + previewBtn + " " + shellBtn + " " + searchBtn + " " + saveBtn + " " + themeBtn + "\n"
	
	// Separator
	separator := t.BorderStyle().Render("") + "\n"
	
	return header + toolbar + separator
}

func (m *model) viewFooter() string {
	t := m.theme.GetTheme()
	
	// Separator
	s := "\n" + t.BorderStyle().Render("") + "\n"
	
	// Left side: status message
	statusMsg := m.core.StatusMessage
	if statusMsg == "" {
		statusMsg = "Ready"
	}
	
	// Right side: file info
	fileInfo := ""
	if m.core.Editor.FilePath != "" {
		fileInfo = filepath.Base(m.core.Editor.FilePath)
	}
	
	// Modified indicator
	modified := ""
	if m.core.Editor.Modified {
		modified = t.ModeStyle("MODIFIED").Render(" + ")
	}
	
	// Build status bar
	statusBar := t.StatusBarStyle().Render(statusMsg) + 
		" " + 
		modified + 
		t.StatusBarStyle().Render(fileInfo) + "\n"
	
	// Help line
	modeStr := m.core.Editor.GetModeString()
	
	help := fmt.Sprintf(" %s | :cmd | /search | \\shell | Tab explorer | Shift+Tab preview | Ctrl+G back | :theme", modeStr)
	statusBar += t.StatusBarStyle().Render(help)
	
	return s + statusBar
}

func (m *model) viewEditor() string {
	t := m.theme.GetTheme()
	e := m.core.Editor
	var lines []string
	
	startLine := m.editorScroll
	endLine := startLine + m.height - m.headerHeight - m.footerHeight - 5
	if endLine > len(e.Lines) {
		endLine = len(e.Lines)
	}
	
	// Ensure startLine is within bounds
	if startLine > len(e.Lines) {
		startLine = len(e.Lines)
		if startLine > 0 {
			startLine--
		}
	}
	
	for i := startLine; i < endLine; i++ {
		if i >= len(e.Lines) {
			break
		}
		
		line := e.Lines[i]
		isCurrentLine := i == e.CursorY
		
		// Line number
		lineNum := fmt.Sprintf(" %3d ", i+1)
		lineNumStyle := t.LineNumberStyle(isCurrentLine)
		
		// Current line indicator
		if isCurrentLine {
			lineNum = "→" + lineNum[1:]
		}
		
		// Line content with syntax highlighting
		var lineContent string
		if isCurrentLine && e.Mode == editor.ModeInsert {
			// Show cursor in insert mode
			before := ""
			after := ""
			cursor := ""
			if e.CursorX < len(line) {
				before = line[:e.CursorX]
				cursor = string(line[e.CursorX])
				after = line[e.CursorX+1:]
			} else {
				before = line
				cursor = " "
			}
			
			lineContent = t.EditorLineStyle(i, isCurrentLine).Render(before) +
				t.CursorStyle().Render(cursor) +
				t.EditorLineStyle(i, isCurrentLine).Render(after)
		} else if e.Selection.Active {
			// Show selection
			lineContent = m.renderLineWithSelection(i, line, isCurrentLine)
		} else {
			// Apply markdown highlighting
			lineContent = m.renderMarkdownLine(i, line, isCurrentLine)
		}
		
		// Combine line number and content
		fullLine := lineNumStyle.Render(lineNum) + lineContent
		lines = append(lines, fullLine)
	}
	
	if len(lines) == 0 {
		lines = []string{lipgloss.NewStyle().Foreground(t.Comment).Render(" (empty file) ")}
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	
	borderHeight := m.height - m.headerHeight - m.footerHeight
	if borderHeight < 5 {
		borderHeight = 10
	}
	
	return t.BorderStyle().Render(content)
}

func (m *model) renderMarkdownLine(lineNum int, line string, isCurrentLine bool) string {
	t := m.theme.GetTheme()
	
	// Simple markdown highlighting
	if strings.HasPrefix(line, "# ") {
		return t.ModeStyle("H1").Render(line)
	} else if strings.HasPrefix(line, "## ") {
		return lipgloss.NewStyle().Foreground(t.Yellow).Bold(true).Render(line)
	} else if strings.HasPrefix(line, "### ") {
		return lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render(line)
	} else if strings.HasPrefix(line, "> ") {
		return lipgloss.NewStyle().Foreground(t.Comment).Italic(true).Render(line)
	} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
		bullet := lipgloss.NewStyle().Foreground(t.Accent).Render("•")
		return t.EditorLineStyle(lineNum, isCurrentLine).Render(" ") + bullet + t.EditorLineStyle(lineNum, isCurrentLine).Render(line[2:])
	} else if strings.HasPrefix(line, "```") {
		return lipgloss.NewStyle().Foreground(t.Muted).Render(line)
	} else if strings.Contains(line, "**") {
		// Bold text
		parts := strings.Split(line, "**")
		var result string
		for i, part := range parts {
			if i%2 == 0 {
				result += t.EditorLineStyle(lineNum, isCurrentLine).Render(part)
			} else {
				result += lipgloss.NewStyle().Foreground(t.Foreground).Bold(true).Render(part)
			}
		}
		return result
	} else if strings.Contains(line, "`") {
		// Inline code
		parts := strings.Split(line, "`")
		var result string
		for i, part := range parts {
			if i%2 == 0 {
				result += t.EditorLineStyle(lineNum, isCurrentLine).Render(part)
			} else {
				result += lipgloss.NewStyle().Background(t.BgMuted).Foreground(t.Accent).Render(" " + part + " ")
			}
		}
		return result
	}
	
	return t.EditorLineStyle(lineNum, isCurrentLine).Render(line)
}

func (m *model) renderLineWithSelection(lineNum int, line string, isCurrentLine bool) string {
	t := m.theme.GetTheme()
	e := m.core.Editor
	
	if !e.Selection.Active {
		return t.EditorLineStyle(lineNum, isCurrentLine).Render(line)
	}
	
	startLine, endLine, startCol, endCol := e.GetSelectionBounds()
	
	if lineNum < startLine || lineNum > endLine {
		return t.EditorLineStyle(lineNum, isCurrentLine).Render(line)
	}
	
	// Calculate selection bounds for this line
	lineStart := 0
	lineEnd := len(line)
	
	if lineNum == startLine {
		lineStart = startCol
	}
	if lineNum == endLine {
		if endCol <= len(line) {
			lineEnd = endCol
		}
	}
	
	// Render with selection highlighting
	before := line[:lineStart]
	selected := line[lineStart:lineEnd]
	after := line[lineEnd:]
	
	return t.EditorLineStyle(lineNum, isCurrentLine).Render(before) +
		t.SelectedStyle().Render(selected) +
		t.EditorLineStyle(lineNum, isCurrentLine).Render(after)
}

func (m *model) viewEditorWithExplorer() string {
	t := m.theme.GetTheme()
	
	// Calculate widths
	explorerWidth := 30
	
	// Render explorer
	explorer := m.viewExplorerContent()
	explorerStyled := t.BorderStyle().Width(explorerWidth).Render(explorer)
	
	// Render editor
	editor := m.viewEditor()
	
	return lipgloss.JoinHorizontal(lipgloss.Top, explorerStyled, " ", editor)
}

func (m *model) viewEditorWithPreview() string {
	t := m.theme.GetTheme()
	
	// Calculate widths
	previewWidth := m.width / 2
	
	// Render editor
	editor := m.viewEditor()
	
	// Render preview
	preview := m.viewPreview()
	previewStyled := t.BorderStyle().Width(previewWidth).Render(preview)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, editor, " ", previewStyled)
}

func (m *model) viewEditorWithExplorerAndPreview() string {
	t := m.theme.GetTheme()
	
	// Calculate widths
	explorerWidth := 25
	previewWidth := m.width / 3
	
	// Render explorer
	explorer := m.viewExplorerContent()
	explorerStyled := t.BorderStyle().Width(explorerWidth).Render(explorer)
	
	// Render editor
	editor := m.viewEditor()
	
	// Render preview
	preview := m.viewPreview()
	previewStyled := t.BorderStyle().Width(previewWidth).Render(preview)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, explorerStyled, " ", editor, " ", previewStyled)
}

func (m *model) viewExplorerContent() string {
	t := m.theme.GetTheme()
	dir := m.core.SpaceDir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	
	entries, _ := os.ReadDir(dir)
	var visibleEntries []os.DirEntry
	for _, e := range entries {
		if len(e.Name()) > 0 && e.Name()[0] != '.' {
			visibleEntries = append(visibleEntries, e)
		}
	}
	
	var lines []string
	
	// Header
	lines = append(lines, t.HeaderStyle().Render("Files"))
	lines = append(lines, "")
	
	for i, entry := range visibleEntries {
		name := entry.Name()
		
		// File icons based on extension
		icon := "  "
		style := t.FileStyle()
		if entry.IsDir() {
			icon = "📁 "
			style = t.DirectoryStyle()
		} else {
			// Determine icon by extension
			ext := strings.ToLower(filepath.Ext(name))
			switch ext {
			case ".md":
				icon = "📝 "
			case ".txt":
				icon = "📄 "
			case ".go":
				icon = "🐹 "
			case ".js", ".ts":
				icon = "📜 "
			case ".json", ".yaml", ".yml", ".toml":
				icon = "⚙️ "
			case ".lua":
				icon = "🌙 "
			case ".db":
				icon = "📊 "
			case ".png", ".jpg", ".jpeg", ".gif":
				icon = "🖼️ "
			default:
				icon = "📄 "
			}
		}
		
		line := fmt.Sprintf(" %s%s", icon, name)
		
		if i == m.selectedFile {
			line = t.SelectedFileStyle().Render(line)
		} else {
			line = style.Render(line)
		}
		
		lines = append(lines, line)
	}
	
	if len(lines) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(t.Comment).Render("  (empty)"))
	}
	
	lines = append(lines, "")
	lines = append(lines, t.ButtonStyle().Render(" + new file"))
	lines = append(lines, t.ButtonStyle().Render(" \\ shell"))
	lines = append(lines, t.ButtonStyle().Render(" : command"))
	
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *model) viewExplorer() string {
	t := m.theme.GetTheme()
	
	content := m.viewExplorerContent()
	
	return t.BorderStyle().Render(content)
}

func (m *model) viewPreview() string {
	t := m.theme.GetTheme()
	
	if m.core.Editor.FilePath == "" || !strings.HasSuffix(strings.ToLower(m.core.Editor.FilePath), ".md") {
		return lipgloss.NewStyle().Foreground(t.Comment).Render("  (no markdown file)")
	}
	
	content := m.core.Editor.String()
	
	// Render markdown to terminal
	rendered := m.markdown.RenderToTerminal(content)
	
	// Apply scroll
	lines := strings.Split(rendered, "\n")
	previewHeight := m.height - m.headerHeight - m.footerHeight - 3
	
	if m.previewScroll >= len(lines) {
		m.previewScroll = len(lines) - 1
	}
	if m.previewScroll < 0 {
		m.previewScroll = 0
	}
	
	endLine := m.previewScroll + previewHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}
	
	visibleLines := lines[m.previewScroll:endLine]
	
	// Add header
	header := t.HeaderStyle().Render("Preview") + "\n"
	
	return header + strings.Join(visibleLines, "\n")
}

// Helper function to extract word at position
func extractWordAt(line string, pos int) string {
	if pos >= len(line) || pos < 0 {
		return ""
	}
	
	// Find word boundaries
	start := pos
	for start > 0 && isWordChar(line[start-1]) {
		start--
	}
	
	end := pos
	for end < len(line) && isWordChar(line[end]) {
		end++
	}
	
	return line[start:end]
}

func isWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

func viewToString(v core.View) string {
	switch v {
	case core.ViewExplorer:
		return "FILES"
	case core.ViewEditor:
		return "EDITOR"
	case core.ViewShell:
		return "SHELL"
	case core.ViewSearch:
		return "SEARCH"
	case core.ViewPreview:
		return "PREVIEW"
	default:
		return "UNKNOWN"
	}
}

// Finder functions

func (m *model) openFinder(finderType string) {
	m.finderMode = true
	m.finderType = finderType
	m.finder.Clear()
	
	dir := m.core.SpaceDir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	
	switch finderType {
	case "file":
		items := finder.ScanDirectoryRecursive(dir, 3, false)
		m.finder.SetItems(items)
		m.core.StatusMessage = "File finder - type to search"
	case "tag":
		// Search for tags
		items := finder.FindFilesWithTag(dir, "")
		if len(items) == 0 {
			// If no tags found, show all files
			items = finder.ScanDirectoryRecursive(dir, 3, false)
		}
		m.finder.SetItems(items)
		m.core.StatusMessage = "Tag finder - type to search"
	case "search":
		// Full text search
		m.core.StatusMessage = "Search - type to find text"
	}
}

func (m *model) handleFinderKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.finderMode = false
		m.core.StatusMessage = ""
	case "enter":
		// Select current item
		item := m.finder.GetSelected()
		if item != nil && !item.IsDir {
			m.core.OpenFile(item.Path)
			m.finderMode = false
			m.core.StatusMessage = fmt.Sprintf("Opened: %s", item.Name)
		} else if item != nil && item.IsDir {
			// Navigate into directory
			m.core.SpaceDir = item.Path
			items := finder.ScanDirectory(item.Path, false)
			m.finder.SetItems(items)
			m.finder.Clear()
		}
	case "backspace":
		m.finder.DeleteChar()
	case "up", "k":
		m.finder.MoveUp()
	case "down", "j":
		m.finder.MoveDown()
	case "ctrl+u":
		m.finder.MovePageUp()
	case "ctrl+d":
		m.finder.MovePageDown()
	case "ctrl+a", "home":
		m.finder.MoveHome()
	case "ctrl+e", "end":
		m.finder.MoveEnd()
	case "ctrl+l":
		m.finder.Clear()
	default:
		// Add character to query
		for _, ch := range msg.Runes {
			if ch >= 32 && ch < 127 { // Printable ASCII
				m.finder.AppendQuery(string(ch))
			}
		}
	}
	
	return m, nil
}

func (m *model) viewFinder() string {
	t := m.theme.GetTheme()
	
	title := "🔍 "
	switch m.finderType {
	case "file":
		title += "File Finder"
	case "tag":
		title += "Tag Finder"
	case "search":
		title += "Search"
	}
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(t.Foreground).
		Bold(true).
		Render(title) + "\n"
	
	// Search input
	query := m.finder.GetQuery()
	searchLine := lipgloss.NewStyle().
		Foreground(t.Accent).
		Render("❯ ") + query + "█\n"
	
	// Results
	items := m.finder.GetFiltered()
	maxItems := m.height - 8
	if maxItems > len(items) {
		maxItems = len(items)
	}
	
	var results []string
	selectedIdx := m.finder.GetSelectedIndex()
	
	for i := 0; i < maxItems; i++ {
		if i >= len(items) {
			break
		}
		
		item := items[i]
		line := ""
		
		// Selection indicator
		if i == selectedIdx {
			line = lipgloss.NewStyle().
				Background(t.Selection).
				Foreground(t.Foreground).
				Bold(true).
				Render(" → ")
		} else {
			line = "   "
		}
		
		// Icon
		icon := "📄 "
		if item.IsDir {
			icon = "📁 "
		} else {
			ext := strings.ToLower(item.Name)
			switch {
			case strings.HasSuffix(ext, ".md"):
				icon = "📝 "
			case strings.HasSuffix(ext, ".go"):
				icon = "🐹 "
			case strings.HasSuffix(ext, ".js") || strings.HasSuffix(ext, ".ts"):
				icon = "📜 "
			case strings.HasSuffix(ext, ".json") || strings.HasSuffix(ext, ".yaml") || strings.HasSuffix(ext, ".yml") || strings.HasSuffix(ext, ".toml"):
				icon = "⚙️ "
			}
		}
		line += icon
		
		// Name with highlighting
		name := item.Name
		if query != "" {
			name = highlightMatch(name, query)
		}
		
		if i == selectedIdx {
			line += lipgloss.NewStyle().
				Background(t.Selection).
				Foreground(t.Foreground).
				Render(name)
		} else {
			line += lipgloss.NewStyle().
				Foreground(t.Foreground).
				Render(name)
		}
		
		// Tag/Path info
		info := ""
		if item.Tag != "" {
			info = lipgloss.NewStyle().
				Foreground(t.Comment).
				Render(" [" + item.Tag + "]")
		} else if item.Path != "" {
			dir := filepath.Dir(item.Path)
			if dir != "." {
				info = lipgloss.NewStyle().
					Foreground(t.Comment).
					Render(" " + dir)
			}
		}
		
		line += info
		results = append(results, line)
	}
	
	if len(results) == 0 {
		results = append(results, lipgloss.NewStyle().
			Foreground(t.Comment).
			Render("  No results"))
	}
	
	// Footer
	footer := "\n" + lipgloss.NewStyle().
		Foreground(t.Comment).
		Render("↑↓ navigate • Enter open • Esc cancel • Ctrl+L clear")
	
	content := header + searchLine + "\n" + strings.Join(results, "\n") + footer
	
	return t.BorderStyle().Width(m.width - 4).Render(content)
}

func highlightMatch(s, query string) string {
	if query == "" {
		return s
	}
	
	lowerS := strings.ToLower(s)
	lowerQ := strings.ToLower(query)
	
	idx := strings.Index(lowerS, lowerQ)
	if idx == -1 {
		return s
	}
	
	before := s[:idx]
	match := s[idx : idx+len(query)]
	after := s[idx+len(query):]
	
	return before + "\033[1;33m" + match + "\033[0m" + after
}

// handleSubstitute handles :s/old/new/g command
func (m *model) handleSubstitute(cmd string) int {
	e := m.core.Editor
	
	// Parse :s/old/new/g or :s/old/new
	var old, newStr string
	global := false
	
	if strings.HasPrefix(cmd, "s/") {
		parts := strings.Split(cmd[2:], "/")
		if len(parts) >= 2 {
			old = parts[0]
			newStr = parts[1]
			if len(parts) >= 3 && parts[2] == "g" {
				global = true
			}
		}
	} else if strings.HasPrefix(cmd, "substitute ") {
		// Parse :substitute old new
		args := strings.TrimSpace(strings.TrimPrefix(cmd, "substitute "))
		parts := strings.Fields(args)
		if len(parts) >= 2 {
			old = parts[0]
			newStr = parts[1]
			if len(parts) >= 3 && parts[2] == "g" {
				global = true
			}
		}
	}
	
	if old == "" {
		return 0
	}
	
	// Check for % for global file substitution
	if global {
		return e.SubstituteAll(old, newStr)
	}
	
	return e.SubstituteLine(old, newStr, global)
}

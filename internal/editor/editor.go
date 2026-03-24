package editor

import (
	"strings"
	"time"
)

// Mode represents the editor mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeSearch
	ModeVisual       // Visual character selection
	ModeVisualLine   // Visual line selection
	ModeVisualBlock  // Visual block selection
)

// Selection represents a text selection
type Selection struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
	Active bool // Whether selection is active
}

// Editor represents the text editor
type Editor struct {
	Content     string
	Lines       []string
	CursorX     int
	CursorY     int
	Mode        Mode
	CommandLine string
	CommandBuf  string
	Modified    bool
	FilePath    string
	Selection   Selection
	
	// Undo/Redo
	history     []EditorState
	historyIdx  int
	maxHistory  int
	
	// Yank register
	YankBuffer  string
	YankIsLine  bool
	
	// Search
	SearchPattern string
	SearchResults []SearchResult
}

// EditorState represents a snapshot of the editor
type EditorState struct {
	Lines     []string
	CursorX   int
	CursorY   int
	Timestamp time.Time
}

// SearchResult represents a search match
type SearchResult struct {
	Line   int
	Column int
	Length int
}

// New creates a new editor
func New() *Editor {
	e := &Editor{
		Content:    "",
		Lines:      []string{""},
		CursorX:    0,
		CursorY:    0,
		Mode:       ModeNormal,
		maxHistory: 100,
	}
	e.saveState()
	return e
}

// saveState saves current state to history
func (e *Editor) saveState() {
	state := EditorState{
		Lines:     make([]string, len(e.Lines)),
		CursorX:   e.CursorX,
		CursorY:   e.CursorY,
		Timestamp: time.Now(),
	}
	copy(state.Lines, e.Lines)
	
	// Remove any future states if we're not at the end
	if e.historyIdx < len(e.history)-1 {
		e.history = e.history[:e.historyIdx+1]
	}
	
	e.history = append(e.history, state)
	if len(e.history) > e.maxHistory {
		e.history = e.history[1:]
	}
	e.historyIdx = len(e.history) - 1
}

// Load loads content into the editor
func (e *Editor) Load(content string) {
	e.Content = content
	e.Lines = strings.Split(content, "\n")
	if len(e.Lines) == 0 {
		e.Lines = []string{""}
	}
	e.Modified = false
	e.CursorX = 0
	e.CursorY = 0
	e.history = nil
	e.historyIdx = 0
	e.saveState()
}

// InsertChar inserts a character at cursor
func (e *Editor) InsertChar(ch rune) {
	if e.CursorY >= len(e.Lines) {
		e.Lines = append(e.Lines, "")
	}
	
	line := e.Lines[e.CursorY]
	if e.CursorX > len(line) {
		e.CursorX = len(line)
	}
	
	e.Lines[e.CursorY] = line[:e.CursorX] + string(ch) + line[e.CursorX:]
	e.CursorX++
	e.Modified = true
}

// DeleteChar deletes character before cursor
func (e *Editor) DeleteChar() {
	if e.CursorY >= len(e.Lines) {
		return
	}
	
	line := e.Lines[e.CursorY]
	if e.CursorX > 0 && e.CursorX <= len(line) {
		e.Lines[e.CursorY] = line[:e.CursorX-1] + line[e.CursorX:]
		e.CursorX--
		e.Modified = true
	}
}

// DeleteCharForward deletes character at cursor
func (e *Editor) DeleteCharForward() {
	if e.CursorY >= len(e.Lines) {
		return
	}
	
	line := e.Lines[e.CursorY]
	if e.CursorX < len(line) {
		e.Lines[e.CursorY] = line[:e.CursorX] + line[e.CursorX+1:]
		e.Modified = true
	}
}

// NewLine creates a new line
func (e *Editor) NewLine() {
	if e.CursorY >= len(e.Lines) {
		e.Lines = append(e.Lines, "")
	}
	
	rest := ""
	if e.CursorX < len(e.Lines[e.CursorY]) {
		rest = e.Lines[e.CursorY][e.CursorX:]
		e.Lines[e.CursorY] = e.Lines[e.CursorY][:e.CursorX]
	}
	
	e.Lines = append(e.Lines[:e.CursorY+1], append([]string{rest}, e.Lines[e.CursorY+1:]...)...)
	e.CursorY++
	e.CursorX = 0
	e.Modified = true
}

// OpenLine inserts a new line below (like Vim's 'o')
func (e *Editor) OpenLine() {
	e.CursorX = len(e.Lines[e.CursorY])
	e.NewLine()
}

// OpenLineAbove inserts a new line above (like Vim's 'O')
func (e *Editor) OpenLineAbove() {
	if e.CursorY == 0 {
		e.Lines = append([]string{""}, e.Lines...)
		e.CursorX = 0
	} else {
		e.CursorY--
		e.CursorX = len(e.Lines[e.CursorY])
		e.NewLine()
		e.CursorY--
	}
	e.Modified = true
}

// DeleteLine deletes the current line
func (e *Editor) DeleteLine() {
	if len(e.Lines) == 1 {
		e.Lines = []string{""}
		e.CursorX = 0
		e.CursorY = 0
	} else {
		e.YankBuffer = e.Lines[e.CursorY]
		e.YankIsLine = true
		e.Lines = append(e.Lines[:e.CursorY], e.Lines[e.CursorY+1:]...)
		if e.CursorY >= len(e.Lines) {
			e.CursorY = len(e.Lines) - 1
		}
	}
	e.CursorX = 0
	e.Modified = true
}

// YankLine yanks (copies) the current line
func (e *Editor) YankLine() {
	if e.CursorY < len(e.Lines) {
		e.YankBuffer = e.Lines[e.CursorY]
		e.YankIsLine = true
	}
}

// YankSelection yanks the current selection
func (e *Editor) YankSelection() {
	if !e.Selection.Active {
		e.YankLine()
		return
	}
	
	startLine, endLine, startCol, endCol := e.GetSelectionBounds()
	
	if startLine == endLine {
		// Single line selection
		line := e.Lines[startLine]
		if endCol > len(line) {
			endCol = len(line)
		}
		e.YankBuffer = line[startCol:endCol]
		e.YankIsLine = false
	} else {
		// Multi-line selection
		var parts []string
		for i := startLine; i <= endLine; i++ {
			line := e.Lines[i]
			if i == startLine {
				parts = append(parts, line[startCol:])
			} else if i == endLine {
				if endCol > len(line) {
					endCol = len(line)
				}
				parts = append(parts, line[:endCol])
			} else {
				parts = append(parts, line)
			}
		}
		e.YankBuffer = strings.Join(parts, "\n")
		e.YankIsLine = true
	}
}

// PasteAfter pastes after cursor
func (e *Editor) PasteAfter() {
	if e.YankBuffer == "" {
		return
	}
	
	if e.YankIsLine {
		// Paste as new line below
		if e.CursorY >= len(e.Lines)-1 {
			e.Lines = append(e.Lines, e.YankBuffer)
		} else {
			e.Lines = append(e.Lines[:e.CursorY+1], append([]string{e.YankBuffer}, e.Lines[e.CursorY+1:]...)...)
		}
		e.CursorY++
		e.CursorX = 0
	} else {
		// Paste at cursor position
		line := e.Lines[e.CursorY]
		if e.CursorX >= len(line) {
			e.Lines[e.CursorY] = line + e.YankBuffer
		} else {
			e.Lines[e.CursorY] = line[:e.CursorX] + e.YankBuffer + line[e.CursorX:]
		}
		e.CursorX += len(e.YankBuffer)
	}
	e.Modified = true
}

// PasteBefore pastes before cursor
func (e *Editor) PasteBefore() {
	if e.YankBuffer == "" {
		return
	}
	
	if e.YankIsLine {
		// Paste as new line above
		if e.CursorY == 0 {
			e.Lines = append([]string{e.YankBuffer}, e.Lines...)
		} else {
			e.Lines = append(e.Lines[:e.CursorY], append([]string{e.YankBuffer}, e.Lines[e.CursorY:]...)...)
		}
		e.CursorX = 0
	} else {
		// Paste at cursor position
		line := e.Lines[e.CursorY]
		if e.CursorX >= len(line) {
			e.Lines[e.CursorY] = line + e.YankBuffer
		} else {
			e.Lines[e.CursorY] = line[:e.CursorX] + e.YankBuffer + line[e.CursorX:]
		}
		e.CursorX += len(e.YankBuffer)
	}
	e.Modified = true
}

// ChangeLine changes (deletes and enters insert mode) the current line
func (e *Editor) ChangeLine() {
	e.YankLine()
	e.Lines[e.CursorY] = ""
	e.CursorX = 0
	e.Mode = ModeInsert
	e.Modified = true
}

// Undo undoes the last change
func (e *Editor) Undo() bool {
	if e.historyIdx <= 0 {
		return false
	}
	
	e.historyIdx--
	state := e.history[e.historyIdx]
	e.Lines = make([]string, len(state.Lines))
	copy(e.Lines, state.Lines)
	e.CursorX = state.CursorX
	e.CursorY = state.CursorY
	e.Modified = true
	return true
}

// Redo redoes the last undone change
func (e *Editor) Redo() bool {
	if e.historyIdx >= len(e.history)-1 {
		return false
	}
	
	e.historyIdx++
	state := e.history[e.historyIdx]
	e.Lines = make([]string, len(state.Lines))
	copy(e.Lines, state.Lines)
	e.CursorX = state.CursorX
	e.CursorY = state.CursorY
	e.Modified = true
	return true
}

// Navigation methods

// MoveLeft moves cursor left
func (e *Editor) MoveLeft() {
	if e.CursorX > 0 {
		e.CursorX--
	} else if e.CursorY > 0 {
		e.CursorY--
		e.CursorX = len(e.Lines[e.CursorY])
	}
}

// MoveRight moves cursor right
func (e *Editor) MoveRight() {
	if e.CursorY < len(e.Lines) && e.CursorX < len(e.Lines[e.CursorY]) {
		e.CursorX++
	} else if e.CursorY < len(e.Lines)-1 {
		e.CursorY++
		e.CursorX = 0
	}
}

// MoveUp moves cursor up
func (e *Editor) MoveUp() {
	if e.CursorY > 0 {
		e.CursorY--
		if e.CursorX > len(e.Lines[e.CursorY]) {
			e.CursorX = len(e.Lines[e.CursorY])
		}
	}
}

// MoveDown moves cursor down
func (e *Editor) MoveDown() {
	if e.CursorY < len(e.Lines)-1 {
		e.CursorY++
		if e.CursorX > len(e.Lines[e.CursorY]) {
			e.CursorX = len(e.Lines[e.CursorY])
		}
	}
}

// MoveWordForward moves to next word start
func (e *Editor) MoveWordForward() {
	line := e.Lines[e.CursorY]
	if e.CursorX >= len(line) {
		if e.CursorY < len(e.Lines)-1 {
			e.CursorY++
			e.CursorX = 0
		}
		return
	}
	
	// Skip current word
	inWord := isWordChar(line[e.CursorX])
	for i := e.CursorX + 1; i < len(line); i++ {
		if isWordChar(line[i]) != inWord {
			e.CursorX = i
			return
		}
	}
	
	// Move to next line
	if e.CursorY < len(e.Lines)-1 {
		e.CursorY++
		e.CursorX = 0
	} else {
		e.CursorX = len(line)
	}
}

// MoveWordBackward moves to previous word start
func (e *Editor) MoveWordBackward() {
	if e.CursorX == 0 {
		if e.CursorY > 0 {
			e.CursorY--
			e.CursorX = len(e.Lines[e.CursorY])
		}
		return
	}
	
	line := e.Lines[e.CursorY]
	
	// Skip whitespace
	i := e.CursorX - 1
	for i > 0 && !isWordChar(line[i]) {
		i--
	}
	
	// Find word start
	for i > 0 && isWordChar(line[i-1]) {
		i--
	}
	
	e.CursorX = i
}

// MoveToEndOfWord moves to end of current/next word
func (e *Editor) MoveToEndOfWord() {
	line := e.Lines[e.CursorY]
	if e.CursorX >= len(line) {
		if e.CursorY < len(e.Lines)-1 {
			e.CursorY++
			e.CursorX = 0
			e.MoveToEndOfWord()
		}
		return
	}
	
	// Find end of current word
	for i := e.CursorX; i < len(line); i++ {
		if !isWordChar(line[i]) {
			if i > e.CursorX {
				e.CursorX = i - 1
				return
			}
			break
		}
	}
	
	// Find next word end
	inWord := false
	for i := e.CursorX + 1; i < len(line); i++ {
		if isWordChar(line[i]) {
			if !inWord {
				inWord = true
			}
		} else if inWord {
			e.CursorX = i - 1
			return
		}
	}
	
	e.CursorX = len(line) - 1
}

// MoveToLineStart moves to start of line (0 in Vim)
func (e *Editor) MoveToLineStart() {
	e.CursorX = 0
}

// MoveToLineEnd moves to end of line ($ in Vim)
func (e *Editor) MoveToLineEnd() {
	if e.CursorY < len(e.Lines) {
		e.CursorX = len(e.Lines[e.CursorY])
	}
}

// MoveToFirstNonBlank moves to first non-blank character
func (e *Editor) MoveToFirstNonBlank() {
	line := e.Lines[e.CursorY]
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			e.CursorX = i
			return
		}
	}
	e.CursorX = 0
}

// DeleteWordForward deletes word forward (dw in Vim)
func (e *Editor) DeleteWordForward() {
	e.saveState()
	startX := e.CursorX
	e.MoveWordForward()
	
	line := e.Lines[e.CursorY]
	if startX < len(line) {
		e.Lines[e.CursorY] = line[:startX] + line[e.CursorX:]
	} else if e.CursorY < len(e.Lines)-1 {
		// Delete line break
		e.Lines[e.CursorY] = line + e.Lines[e.CursorY+1]
		e.Lines = append(e.Lines[:e.CursorY+1], e.Lines[e.CursorY+2:]...)
	}
	
	e.CursorX = startX
	e.Modified = true
}

// ChangeWordForward changes word forward (cw in Vim)
func (e *Editor) ChangeWordForward() {
	e.saveState()
	startX := e.CursorX
	e.MoveWordForward()
	
	line := e.Lines[e.CursorY]
	if startX < len(line) {
		e.Lines[e.CursorY] = line[:startX] + line[e.CursorX:]
	}
	
	e.CursorX = startX
	e.Mode = ModeInsert
	e.Modified = true
}

// JoinLines joins current line with next (J in Vim)
func (e *Editor) JoinLines() {
	if e.CursorY >= len(e.Lines)-1 {
		return
	}
	
	e.saveState()
	line := e.Lines[e.CursorY]
	nextLine := strings.TrimSpace(e.Lines[e.CursorY+1])
	
	if len(line) > 0 && line[len(line)-1] != ' ' {
		line += " "
	}
	
	e.Lines[e.CursorY] = line + nextLine
	e.Lines = append(e.Lines[:e.CursorY+1], e.Lines[e.CursorY+2:]...)
	e.Modified = true
}

// IndentLine indents the current line (>> in Vim)
func (e *Editor) IndentLine() {
	e.saveState()
	e.Lines[e.CursorY] = "  " + e.Lines[e.CursorY]
	e.CursorX += 2
	e.Modified = true
}

// OutdentLine removes indent from current line (<< in Vim)
func (e *Editor) OutdentLine() {
	e.saveState()
	line := e.Lines[e.CursorY]
	if strings.HasPrefix(line, "  ") {
		e.Lines[e.CursorY] = line[2:]
		if e.CursorX >= 2 {
			e.CursorX -= 2
		} else {
			e.CursorX = 0
		}
		e.Modified = true
	}
}

// ToggleCase toggles case of character under cursor (~ in Vim)
func (e *Editor) ToggleCase() {
	if e.CursorY >= len(e.Lines) {
		return
	}
	
	line := e.Lines[e.CursorY]
	if e.CursorX >= len(line) {
		return
	}
	
	ch := line[e.CursorX]
	if ch >= 'a' && ch <= 'z' {
		ch = ch - 'a' + 'A'
	} else if ch >= 'A' && ch <= 'Z' {
		ch = ch - 'A' + 'a'
	}
	
	e.Lines[e.CursorY] = line[:e.CursorX] + string(ch) + line[e.CursorX+1:]
	e.Modified = true
}

// Visual Mode methods

// EnterVisualMode enters character-wise visual mode
func (e *Editor) EnterVisualMode() {
	e.Mode = ModeVisual
	e.Selection = Selection{
		StartX: e.CursorX,
		StartY: e.CursorY,
		EndX:   e.CursorX,
		EndY:   e.CursorY,
		Active: true,
	}
}

// EnterVisualLineMode enters line-wise visual mode
func (e *Editor) EnterVisualLineMode() {
	e.Mode = ModeVisualLine
	e.Selection = Selection{
		StartX: 0,
		StartY: e.CursorY,
		EndX:   len(e.Lines[e.CursorY]),
		EndY:   e.CursorY,
		Active: true,
	}
}

// ExitVisualMode exits visual mode
func (e *Editor) ExitVisualMode() {
	e.Mode = ModeNormal
	e.Selection.Active = false
}

// UpdateVisualSelection updates the selection based on cursor position
func (e *Editor) UpdateVisualSelection() {
	e.Selection.EndX = e.CursorX
	e.Selection.EndY = e.CursorY
}

// GetSelectedText returns the selected text
func (e *Editor) GetSelectedText() string {
	if !e.Selection.Active {
		return ""
	}
	
	startLine, endLine, startCol, endCol := e.GetSelectionBounds()
	
	if startLine == endLine {
		line := e.Lines[startLine]
		if endCol > len(line) {
			endCol = len(line)
		}
		return line[startCol:endCol]
	}
	
	var parts []string
	for i := startLine; i <= endLine; i++ {
		line := e.Lines[i]
		if i == startLine {
			parts = append(parts, line[startCol:])
		} else if i == endLine {
			if endCol > len(line) {
				endCol = len(line)
			}
			parts = append(parts, line[:endCol])
		} else {
			parts = append(parts, line)
		}
	}
	return strings.Join(parts, "\n")
}

// GetSelectionBounds returns normalized selection bounds
func (e *Editor) GetSelectionBounds() (startLine, endLine, startCol, endCol int) {
	sel := e.Selection
	
	if sel.StartY < sel.EndY || (sel.StartY == sel.EndY && sel.StartX <= sel.EndX) {
		startLine, endLine = sel.StartY, sel.EndY
		startCol, endCol = sel.StartX, sel.EndX
	} else {
		startLine, endLine = sel.EndY, sel.StartY
		startCol, endCol = sel.EndX, sel.StartX
	}
	
	// Adjust for inclusive selection
	if startLine == endLine {
		endCol++
	} else {
		endCol++
	}
	
	return
}

// DeleteSelection deletes the selected text
func (e *Editor) DeleteSelection() {
	if !e.Selection.Active {
		return
	}
	
	e.saveState()
	startLine, endLine, startCol, endCol := e.GetSelectionBounds()
	
	if startLine == endLine {
		// Single line
		line := e.Lines[startLine]
		if endCol > len(line) {
			endCol = len(line)
		}
		e.Lines[startLine] = line[:startCol] + line[endCol:]
	} else {
		// Multi-line
		firstLine := e.Lines[startLine][:startCol]
		lastLine := ""
		if endCol <= len(e.Lines[endLine]) {
			lastLine = e.Lines[endLine][endCol:]
		}
		e.Lines[startLine] = firstLine + lastLine
		e.Lines = append(e.Lines[:startLine+1], e.Lines[endLine+1:]...)
	}
	
	e.CursorX = startCol
	e.CursorY = startLine
	e.ExitVisualMode()
	e.Modified = true
}

// Search methods

// Search forward for pattern
func (e *Editor) Search(pattern string) bool {
	e.SearchPattern = pattern
	if pattern == "" {
		return false
	}
	
	startX := e.CursorX + 1
	startY := e.CursorY
	
	// Search from cursor to end
	for y := startY; y < len(e.Lines); y++ {
		line := e.Lines[y]
		startIdx := 0
		if y == startY {
			startIdx = startX
		}
		
		idx := strings.Index(line[startIdx:], pattern)
		if idx != -1 {
			e.CursorY = y
			e.CursorX = startIdx + idx
			return true
		}
	}
	
	// Wrap around to beginning
	for y := 0; y <= startY; y++ {
		line := e.Lines[y]
		endIdx := len(line)
		if y == startY {
			endIdx = startX
		}
		
		idx := strings.Index(line[:endIdx], pattern)
		if idx != -1 {
			e.CursorY = y
			e.CursorX = idx
			return true
		}
	}
	
	return false
}

// SearchNext finds next occurrence
func (e *Editor) SearchNext() bool {
	return e.Search(e.SearchPattern)
}

// SearchPrevious finds previous occurrence
func (e *Editor) SearchPrevious() bool {
	if e.SearchPattern == "" {
		return false
	}
	
	// Simple implementation - search backward
	for y := e.CursorY; y >= 0; y-- {
		line := e.Lines[y]
		endIdx := e.CursorX
		if y != e.CursorY {
			endIdx = len(line)
		}
		
		idx := strings.LastIndex(line[:endIdx], e.SearchPattern)
		if idx != -1 {
			e.CursorY = y
			e.CursorX = idx
			return true
		}
	}
	
	return false
}

// Helper functions

func isWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

// Command mode methods

// EnterCommandMode enters command mode
func (e *Editor) EnterCommandMode() {
	e.Mode = ModeCommand
	e.CommandLine = ""
}

// ExitCommandMode exits command mode
func (e *Editor) ExitCommandMode() {
	e.Mode = ModeNormal
	e.CommandLine = ""
}

// EnterSearchMode enters search mode
func (e *Editor) EnterSearchMode() {
	e.Mode = ModeSearch
	e.CommandBuf = ""
}

// ExitSearchMode exits search mode
func (e *Editor) ExitSearchMode() {
	e.Mode = ModeNormal
	e.CommandBuf = ""
}

// TypeCommand types a character in command mode
func (e *Editor) TypeCommand(ch rune) {
	if ch == 27 { // Escape
		e.ExitCommandMode()
		return
	}
	e.CommandLine += string(ch)
}

// TypeSearch types a character in search mode
func (e *Editor) TypeSearch(ch rune) {
	if ch == 27 { // Escape
		e.ExitSearchMode()
		return
	}
	e.CommandBuf += string(ch)
}

// ExecuteCommand executes a command line command
func (e *Editor) ExecuteCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	
	command := strings.ToLower(parts[0])
	args := parts[1:]
	
	switch command {
	case "wq":
		return "quit_save"
	case "q":
		return "quit"
	case "q!":
		return "quit_force"
	case "w":
		return "write"
	case "w!":
		return "write_force"
	case "e", "edit":
		if len(args) > 0 {
			return "edit:" + args[0]
		}
	case "new":
		return "new_file"
	case "split":
		return "split"
	case "vsplit":
		return "vsplit"
	case "set":
		if len(args) > 0 {
			return "set:" + strings.Join(args, " ")
		}
	default:
		// KV command: key = value
		if strings.Contains(cmd, "=") {
			return "kv:" + cmd
		}
	}
	
	return ""
}

// String returns the editor content as string
func (e *Editor) String() string {
	return strings.Join(e.Lines, "\n")
}

// GetCurrentLine returns the current line
func (e *Editor) GetCurrentLine() string {
	if e.CursorY >= 0 && e.CursorY < len(e.Lines) {
		return e.Lines[e.CursorY]
	}
	return ""
}

// GetModeString returns a string representation of the current mode
func (e *Editor) GetModeString() string {
	switch e.Mode {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeCommand:
		return "COMMAND"
	case ModeSearch:
		return "SEARCH"
	case ModeVisual:
		return "VISUAL"
	case ModeVisualLine:
		return "V-LINE"
	case ModeVisualBlock:
		return "V-BLOCK"
	default:
		return "UNKNOWN"
	}
}

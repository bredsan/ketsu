package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines color scheme for the UI
type Theme struct {
	// Background colors
	BgDark   lipgloss.Color
	BgLight  lipgloss.Color
	BgMuted  lipgloss.Color
	
	// Text colors
	Foreground lipgloss.Color
	Comment    lipgloss.Color
	Muted      lipgloss.Color
	
	// Accent colors
	Accent  lipgloss.Color
	Green   lipgloss.Color
	Red     lipgloss.Color
	Yellow  lipgloss.Color
	Purple  lipgloss.Color
	Blue    lipgloss.Color
	Orange  lipgloss.Color
	Cyan    lipgloss.Color
	
	// UI colors
	Border     lipgloss.Color
	Selection  lipgloss.Color
	Highlight  lipgloss.Color
	StatusBar  lipgloss.Color
	StatusBarBg lipgloss.Color
	
	// Mode colors
	NormalMode  lipgloss.Color
	InsertMode  lipgloss.Color
	CommandMode lipgloss.Color
	VisualMode  lipgloss.Color
	ReplaceMode lipgloss.Color
}

// Catppuccin Mocha theme
var CatppuccinTheme = Theme{
	BgDark:   lipgloss.Color("#1e1e2e"),
	BgLight:  lipgloss.Color("#302d41"),
	BgMuted:  lipgloss.Color("#45475a"),
	
	Foreground: lipgloss.Color("#cdd6f4"),
	Comment:    lipgloss.Color("#6c7086"),
	Muted:      lipgloss.Color("#a6adc8"),
	
	Accent: lipgloss.Color("#89b4fa"),
	Green:  lipgloss.Color("#a6e3a1"),
	Red:    lipgloss.Color("#f38ba8"),
	Yellow: lipgloss.Color("#f9e2af"),
	Purple: lipgloss.Color("#cba6f7"),
	Blue:   lipgloss.Color("#89b4fa"),
	Orange: lipgloss.Color("#fab387"),
	Cyan:   lipgloss.Color("#89dceb"),
	
	Border:      lipgloss.Color("#45475a"),
	Selection:   lipgloss.Color("#585b70"),
	Highlight:   lipgloss.Color("#7f849c"),
	StatusBar:   lipgloss.Color("#a6adc8"),
	StatusBarBg: lipgloss.Color("#181825"),
	
	NormalMode:  lipgloss.Color("#a6e3a1"),
	InsertMode:  lipgloss.Color("#89b4fa"),
	CommandMode: lipgloss.Color("#f9e2af"),
	VisualMode:  lipgloss.Color("#cba6f7"),
	ReplaceMode: lipgloss.Color("#fab387"),
}

// Tokyo Night theme
var TokyoNightTheme = Theme{
	BgDark:   lipgloss.Color("#1a1b26"),
	BgLight:  lipgloss.Color("#24283b"),
	BgMuted:  lipgloss.Color("#414868"),
	
	Foreground: lipgloss.Color("#c0caf5"),
	Comment:    lipgloss.Color("#565f89"),
	Muted:      lipgloss.Color("#a9b1d6"),
	
	Accent: lipgloss.Color("#7aa2f7"),
	Green:  lipgloss.Color("#9ece6a"),
	Red:    lipgloss.Color("#f7768e"),
	Yellow: lipgloss.Color("#e0af68"),
	Purple: lipgloss.Color("#bb9af7"),
	Blue:   lipgloss.Color("#7aa2f7"),
	Orange: lipgloss.Color("#ff9e64"),
	Cyan:   lipgloss.Color("#7dcfff"),
	
	Border:      lipgloss.Color("#414868"),
	Selection:   lipgloss.Color("#3b4261"),
	Highlight:   lipgloss.Color("#565f89"),
	StatusBar:   lipgloss.Color("#a9b1d6"),
	StatusBarBg: lipgloss.Color("#16161e"),
	
	NormalMode:  lipgloss.Color("#9ece6a"),
	InsertMode:  lipgloss.Color("#7aa2f7"),
	CommandMode: lipgloss.Color("#e0af68"),
	VisualMode:  lipgloss.Color("#bb9af7"),
	ReplaceMode: lipgloss.Color("#ff9e64"),
}

// Nord theme
var NordTheme = Theme{
	BgDark:   lipgloss.Color("#2e3440"),
	BgLight:  lipgloss.Color("#3b4252"),
	BgMuted:  lipgloss.Color("#434c5e"),
	
	Foreground: lipgloss.Color("#d8dee9"),
	Comment:    lipgloss.Color("#4c566a"),
	Muted:      lipgloss.Color("#d8dee9"),
	
	Accent: lipgloss.Color("#88c0d0"),
	Green:  lipgloss.Color("#a3be8c"),
	Red:    lipgloss.Color("#bf616a"),
	Yellow: lipgloss.Color("#ebcb8b"),
	Purple: lipgloss.Color("#b48ead"),
	Blue:   lipgloss.Color("#81a1c1"),
	Orange: lipgloss.Color("#d08770"),
	Cyan:   lipgloss.Color("#8fbcbb"),
	
	Border:      lipgloss.Color("#4c566a"),
	Selection:   lipgloss.Color("#434c5e"),
	Highlight:   lipgloss.Color("#4c566a"),
	StatusBar:   lipgloss.Color("#d8dee9"),
	StatusBarBg: lipgloss.Color("#2e3440"),
	
	NormalMode:  lipgloss.Color("#a3be8c"),
	InsertMode:  lipgloss.Color("#81a1c1"),
	CommandMode: lipgloss.Color("#ebcb8b"),
	VisualMode:  lipgloss.Color("#b48ead"),
	ReplaceMode: lipgloss.Color("#d08770"),
}

// ThemeManager manages the current theme
type ThemeManager struct {
	current Theme
	name    string
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		current: CatppuccinTheme,
		name:    "catppuccin",
	}
}

// SetTheme sets the current theme by name
func (tm *ThemeManager) SetTheme(name string) bool {
	switch name {
	case "catppuccin", "catppuccin-mocha":
		tm.current = CatppuccinTheme
		tm.name = "catppuccin"
		return true
	case "tokyo-night", "tokyonight":
		tm.current = TokyoNightTheme
		tm.name = "tokyo-night"
		return true
	case "nord":
		tm.current = NordTheme
		tm.name = "nord"
		return true
	default:
		return false
	}
}

// GetTheme returns the current theme
func (tm *ThemeManager) GetTheme() Theme {
	return tm.current
}

// GetName returns the current theme name
func (tm *ThemeManager) GetName() string {
	return tm.name
}

// GetAvailableThemes returns list of available theme names
func (tm *ThemeManager) GetAvailableThemes() []string {
	return []string{"catppuccin", "tokyo-night", "nord"}
}

// Style helpers for common UI elements
func (t Theme) HeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Foreground).
		Bold(true).
		Padding(0, 1)
}

func (t Theme) ButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Foreground).
		Background(t.BgMuted).
		Padding(0, 1).
		MarginRight(1)
}

func (t Theme) ButtonActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.BgDark).
		Background(t.Accent).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)
}

func (t Theme) BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
}

func (t Theme) SelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Selection).
		Foreground(t.Foreground)
}

func (t Theme) StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.StatusBar).
		Background(t.StatusBarBg).
		Padding(0, 1)
}

func (t Theme) ModeStyle(mode string) lipgloss.Style {
	var color lipgloss.Color
	switch mode {
	case "INSERT":
		color = t.InsertMode
	case "COMMAND":
		color = t.CommandMode
	case "VISUAL":
		color = t.VisualMode
	case "REPLACE":
		color = t.ReplaceMode
	default:
		color = t.NormalMode
	}
	
	return lipgloss.NewStyle().
		Foreground(t.BgDark).
		Background(color).
		Bold(true).
		Padding(0, 1)
}

func (t Theme) EditorLineStyle(lineNum int, isCurrentLine bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(t.Foreground)
	
	if isCurrentLine {
		style = style.
			Background(t.Selection).
			Bold(true)
	}
	
	return style
}

func (t Theme) LineNumberStyle(isCurrentLine bool) lipgloss.Style {
	if isCurrentLine {
		return lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true)
	}
	return lipgloss.NewStyle().
		Foreground(t.Comment)
}

func (t Theme) ExplorerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Foreground)
}

func (t Theme) DirectoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)
}

func (t Theme) FileStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Foreground)
}

func (t Theme) SelectedFileStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Selection).
		Foreground(t.Foreground).
		Bold(true)
}

func (t Theme) CursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Accent).
		Foreground(t.BgDark)
}

func (t Theme) SelectionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Selection)
}
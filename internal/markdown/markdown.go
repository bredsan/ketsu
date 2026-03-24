package markdown

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Renderer handles Markdown to HTML conversion
type Renderer struct {
	md goldmark.Markdown
}

// New creates a new Markdown renderer
func New() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &Renderer{md: md}
}

// Render converts Markdown to HTML
func (r *Renderer) Render(source string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderToTerminal converts Markdown to styled terminal output
func (r *Renderer) RenderToTerminal(source string) string {
	lines := strings.Split(source, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimRight(line, " \t")
		
		if strings.HasPrefix(line, "# ") {
			// H1 - Bold, large
			result = append(result, "\033[1;36m"+strings.TrimPrefix(line, "# ")+"\033[0m")
		} else if strings.HasPrefix(line, "## ") {
			// H2 - Bold
			result = append(result, "\033[1;33m"+strings.TrimPrefix(line, "## ")+"\033[0m")
		} else if strings.HasPrefix(line, "### ") {
			// H3
			result = append(result, "\033[1;32m"+strings.TrimPrefix(line, "### ")+"\033[0m")
		} else if strings.HasPrefix(line, "#### ") {
			// H4
			result = append(result, "\033[1;35m"+strings.TrimPrefix(line, "#### ")+"\033[0m")
		} else if strings.HasPrefix(line, "> ") {
			// Blockquote
			result = append(result, "\033[3m\033[90m"+line+"\033[0m")
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			// List item
			result = append(result, "\033[33m•\033[0m "+line[2:])
		} else if strings.HasPrefix(line, "```") {
			// Code block marker
			result = append(result, "\033[90m"+line+"\033[0m")
		} else if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			// Indented code
			result = append(result, "\033[90m"+line+"\033[0m")
		} else if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "***") {
			// Horizontal rule
			result = append(result, "\033[90m────────────────────────────────\033[0m")
		} else if line == "" {
			result = append(result, "")
		} else {
			// Process inline formatting
			processed := processInline(line)
			result = append(result, processed)
		}
	}

	return strings.Join(result, "\n")
}

// processInline handles inline Markdown formatting
func processInline(line string) string {
	// Bold **text** or __text__
	line = processBold(line)
	// Italic *text* or _text_
	line = processItalic(line)
	// Inline code `text`
	line = processCode(line)
	// Links [text](url)
	line = processLinks(line)
	
	return line
}

func processBold(line string) string {
	// Handle **text**
	for {
		start := strings.Index(line, "**")
		if start == -1 {
			break
		}
		end := strings.Index(line[start+2:], "**")
		if end == -1 {
			break
		}
		end += start + 2
		content := line[start+2 : end]
		line = line[:start] + "\033[1m" + content + "\033[0m" + line[end+2:]
	}
	
	// Handle __text__
	for {
		start := strings.Index(line, "__")
		if start == -1 {
			break
		}
		end := strings.Index(line[start+2:], "__")
		if end == -1 {
			break
		}
		end += start + 2
		content := line[start+2 : end]
		line = line[:start] + "\033[1m" + content + "\033[0m" + line[end+2:]
	}
	
	return line
}

func processItalic(line string) string {
	// Handle *text* (but not ** which is bold)
	for {
		start := strings.Index(line, "*")
		if start == -1 {
			break
		}
		// Skip if it's part of **
		if start > 0 && line[start-1] == '*' {
			continue
		}
		if start < len(line)-1 && line[start+1] == '*' {
			continue
		}
		end := strings.Index(line[start+1:], "*")
		if end == -1 {
			break
		}
		end += start + 1
		content := line[start+1 : end]
		line = line[:start] + "\033[3m" + content + "\033[0m" + line[end+1:]
	}
	
	return line
}

func processCode(line string) string {
	// Handle `text`
	for {
		start := strings.Index(line, "`")
		if start == -1 {
			break
		}
		end := strings.Index(line[start+1:], "`")
		if end == -1 {
			break
		}
		end += start + 1
		content := line[start+1 : end]
		line = line[:start] + "\033[90m\033[47m " + content + " \033[0m" + line[end+1:]
	}
	
	return line
}

func processLinks(line string) string {
	// Handle [text](url)
	for {
		start := strings.Index(line, "[")
		if start == -1 {
			break
		}
		middle := strings.Index(line[start:], "](")
		if middle == -1 {
			break
		}
		middle += start
		end := strings.Index(line[middle+2:], ")")
		if end == -1 {
			break
		}
		end += middle + 2
		
		text := line[start+1 : middle]
		url := line[middle+2 : end]
		line = line[:start] + "\033[4;34m" + text + "\033[0m" + " \033[90m(" + url + ")\033[0m" + line[end+1:]
	}
	
	return line
}

// ExtractTitle extracts the first heading from Markdown
func ExtractTitle(source string) string {
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return "Untitled"
}

// GetWordCount returns the word count of the text
func GetWordCount(source string) int {
	words := strings.Fields(source)
	return len(words)
}

// GetLineCount returns the line count
func GetLineCount(source string) int {
	return len(strings.Split(source, "\n"))
}

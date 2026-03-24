package finder

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
)

// Item represents a selectable item in the finder
type Item struct {
	Name    string
	Path    string
	IsDir   bool
	Score   int
	Tag     string // For tag-based search
	Line    int    // For search results
	Content string // Preview content
}

// Finder provides fuzzy finding functionality
type Finder struct {
	items     []Item
	filtered  []Item
	query     string
	selected  int
	maxHeight int
}

// New creates a new Finder
func New() *Finder {
	return &Finder{
		items:     make([]Item, 0),
		filtered:  make([]Item, 0),
		selected:  0,
		maxHeight: 15,
	}
}

// SetItems sets the items to search through
func (f *Finder) SetItems(items []Item) {
	f.items = items
	f.filtered = items
	f.selected = 0
}

// SetQuery sets the search query
func (f *Finder) SetQuery(query string) {
	f.query = query
	f.filter()
	f.selected = 0
}

// GetQuery returns the current query
func (f *Finder) GetQuery() string {
	return f.query
}

// GetFiltered returns the filtered items
func (f *Finder) GetFiltered() []Item {
	return f.filtered
}

// GetSelected returns the currently selected item
func (f *Finder) GetSelected() *Item {
	if len(f.filtered) == 0 || f.selected >= len(f.filtered) {
		return nil
	}
	return &f.filtered[f.selected]
}

// GetSelectedIndex returns the selected index
func (f *Finder) GetSelectedIndex() int {
	return f.selected
}

// MoveUp moves selection up
func (f *Finder) MoveUp() {
	if f.selected > 0 {
		f.selected--
	}
}

// MoveDown moves selection down
func (f *Finder) MoveDown() {
	if f.selected < len(f.filtered)-1 {
		f.selected++
	}
}

// MovePageUp moves selection up by page
func (f *Finder) MovePageUp() {
	f.selected -= 5
	if f.selected < 0 {
		f.selected = 0
	}
}

// MovePageDown moves selection down by page
func (f *Finder) MovePageDown() {
	f.selected += 5
	if f.selected >= len(f.filtered) {
		f.selected = len(f.filtered) - 1
	}
}

// MoveHome moves selection to start
func (f *Finder) MoveHome() {
	f.selected = 0
}

// MoveEnd moves selection to end
func (f *Finder) MoveEnd() {
	if len(f.filtered) > 0 {
		f.selected = len(f.filtered) - 1
	}
}

// AppendQuery appends to the query
func (f *Finder) AppendQuery(s string) {
	f.query += s
	f.filter()
	f.selected = 0
}

// DeleteChar deletes last character from query
func (f *Finder) DeleteChar() {
	if len(f.query) > 0 {
		f.query = f.query[:len(f.query)-1]
		f.filter()
		f.selected = 0
	}
}

// Clear clears the query
func (f *Finder) Clear() {
	f.query = ""
	f.filtered = f.items
	f.selected = 0
}

// Len returns the number of filtered items
func (f *Finder) Len() int {
	return len(f.filtered)
}

// filter filters items based on query
func (f *Finder) filter() {
	if f.query == "" {
		f.filtered = f.items
		return
	}

	names := make([]string, len(f.items))
	for i, item := range f.items {
		names[i] = item.Name
		if item.Tag != "" {
			names[i] += " " + item.Tag
		}
		if item.Content != "" {
			names[i] += " " + item.Content
		}
	}

	matches := fuzzy.Find(f.query, names)
	
	f.filtered = make([]Item, 0, len(matches))
	for _, match := range matches {
		item := f.items[match.Index]
		item.Score = match.Score
		f.filtered = append(f.filtered, item)
	}

	// Sort by score (higher is better)
	sort.Slice(f.filtered, func(i, j int) bool {
		return f.filtered[i].Score > f.filtered[j].Score
	})
}

// ScanDirectory scans a directory for files
func ScanDirectory(dir string, includeHidden bool) []Item {
	items := make([]Item, 0)
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return items
	}
	
	for _, entry := range entries {
		name := entry.Name()
		
		// Skip hidden files unless requested
		if !includeHidden && strings.HasPrefix(name, ".") {
			continue
		}
		
		path := filepath.Join(dir, name)
		item := Item{
			Name:  name,
			Path:  path,
			IsDir: entry.IsDir(),
		}
		
		// Add file extension as tag
		if !entry.IsDir() {
			ext := filepath.Ext(name)
			if ext != "" {
				item.Tag = ext[1:] // Remove the dot
			}
		}
		
		items = append(items, item)
	}
	
	// Sort: directories first, then alphabetically
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
	
	return items
}

// ScanDirectoryRecursive scans a directory recursively
func ScanDirectoryRecursive(dir string, maxDepth int, includeHidden bool) []Item {
	items := make([]Item, 0)
	scanDir(dir, "", &items, 0, maxDepth, includeHidden)
	return items
}

func scanDir(baseDir, relPath string, items *[]Item, depth, maxDepth int, includeHidden bool) {
	if depth > maxDepth {
		return
	}
	
	fullPath := filepath.Join(baseDir, relPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return
	}
	
	for _, entry := range entries {
		name := entry.Name()
		
		if !includeHidden && strings.HasPrefix(name, ".") {
			continue
		}
		
		itemRelPath := relPath
		if itemRelPath != "" {
			itemRelPath = filepath.Join(itemRelPath, name)
		} else {
			itemRelPath = name
		}
		
		item := Item{
			Name:  name,
			Path:  filepath.Join(baseDir, itemRelPath),
			IsDir: entry.IsDir(),
		}
		
		if !entry.IsDir() {
			ext := filepath.Ext(name)
			if ext != "" {
				item.Tag = ext[1:]
			}
		} else {
			item.Tag = "dir"
		}
		
		*items = append(*items, item)
		
		// Recurse into subdirectories
		if entry.IsDir() {
			scanDir(baseDir, itemRelPath, items, depth+1, maxDepth, includeHidden)
		}
	}
}

// ExtractTags extracts tags (words starting with #) from content
func ExtractTags(content string) []string {
	tags := make(map[string]bool)
	words := strings.Fields(content)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			// Remove punctuation
			tag = strings.Trim(tag, ".,;:!?()[]{}\"'")
			if tag != "" {
				tags[tag] = true
			}
		}
	}
	
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	
	sort.Strings(result)
	return result
}

// FindFilesWithTag finds files containing a specific tag
func FindFilesWithTag(dir string, tag string) []Item {
	items := make([]Item, 0)
	
	files := ScanDirectoryRecursive(dir, 10, false)
	for _, file := range files {
		if file.IsDir {
			continue
		}
		
		// Only check markdown files
		if !strings.HasSuffix(strings.ToLower(file.Name), ".md") {
			continue
		}
		
		content, err := os.ReadFile(file.Path)
		if err != nil {
			continue
		}
		
		contentStr := string(content)
		tags := ExtractTags(contentStr)
		
		for _, t := range tags {
			if strings.EqualFold(t, tag) {
				file.Tag = strings.Join(tags, ", ")
				file.Content = getFirstLines(contentStr, 2)
				items = append(items, file)
				break
			}
		}
	}
	
	return items
}

// SearchInFiles searches for text in files
func SearchInFiles(dir string, query string) []Item {
	items := make([]Item, 0)
	
	files := ScanDirectoryRecursive(dir, 10, false)
	for _, file := range files {
		if file.IsDir {
			continue
		}
		
		content, err := os.ReadFile(file.Path)
		if err != nil {
			continue
		}
		
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
				item := Item{
					Name:    file.Name,
					Path:    file.Path,
					IsDir:   false,
					Line:    i + 1,
					Content: strings.TrimSpace(line),
					Tag:     "match",
				}
				items = append(items, item)
			}
		}
	}
	
	return items
}

// getFirstLines returns the first n lines of content
func getFirstLines(content string, n int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, " | ")
}

// Render renders the finder UI
func (f *Finder) Render(width, height int) string {
	var result strings.Builder
	
	// Header
	result.WriteString("🔍 Fuzzy Finder")
	result.WriteString("\n")
	
	// Search input
	result.WriteString("❯ ")
	result.WriteString(f.query)
	result.WriteString("█")
	result.WriteString("\n")
	
	// Results
	maxItems := f.maxHeight
	if maxItems > height-5 {
		maxItems = height - 5
	}
	if maxItems > len(f.filtered) {
		maxItems = len(f.filtered)
	}
	
	startIdx := 0
	if f.selected >= maxItems {
		startIdx = f.selected - maxItems + 1
	}
	
	for i := 0; i < maxItems; i++ {
		idx := startIdx + i
		if idx >= len(f.filtered) {
			break
		}
		
		item := f.filtered[idx]
		
		if idx == f.selected {
			result.WriteString("→ ")
		} else {
			result.WriteString("  ")
		}
		
		// Icon
		if item.IsDir {
			result.WriteString("📁 ")
		} else {
			ext := filepath.Ext(item.Name)
			switch strings.ToLower(ext) {
			case ".md":
				result.WriteString("📝 ")
			case ".go":
				result.WriteString("🐹 ")
			case ".js", ".ts":
				result.WriteString("📜 ")
			case ".json", ".yaml", ".yml", ".toml":
				result.WriteString("⚙️ ")
			default:
				result.WriteString("📄 ")
			}
		}
		
		// Name
		result.WriteString(item.Name)
		
		// Tag
		if item.Tag != "" {
			result.WriteString(" [")
			result.WriteString(item.Tag)
			result.WriteString("]")
		}
		
		// Content preview (for search results)
		if item.Content != "" && item.Line > 0 {
			result.WriteString(" :")
			result.WriteString(strings.TrimSpace(item.Content)[:50])
		}
		
		result.WriteString("\n")
	}
	
	// Footer
	result.WriteString("\n")
	result.WriteString("Type to search • Enter to select • Esc to cancel")
	
	return result.String()
}

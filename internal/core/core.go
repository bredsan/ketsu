package core

import (
	"os"
	"path/filepath"

	"github.com/ketsu/ketsu/internal/kv"
	"github.com/ketsu/ketsu/internal/editor"
)

type View int

const (
	ViewExplorer View = iota
	ViewEditor
	ViewShell
	ViewSearch
	ViewPreview
)

type Core struct {
	CurrentView  View
	Editor       *editor.Editor
	KV           *kv.KV
	SpaceDir     string
	SearchQuery  string
	StatusMessage string
}

func New() *Core {
	return &Core{
		CurrentView: ViewExplorer,
		Editor:      editor.New(),
		KV:          kv.New(),
	}
}

func (c *Core) LoadSpace(dir string) error {
	c.SpaceDir = dir
	
	// Create .ketsu directory if not exists
	ketsuDir := filepath.Join(dir, ".ketsu")
	if err := os.MkdirAll(ketsuDir, 0755); err != nil {
		return err
	}
	
	// Load KV data
	if err := c.KV.Load(dir); err != nil {
		return err
	}
	
	return nil
}

func (c *Core) ProcessCommand(cmd string) string {
	// Handle KV commands: key = value
	if idx := findIndexByte(cmd, '='); idx != -1 {
		key := trimSpace(cmd[:idx])
		value := trimSpace(cmd[idx+1:])
		
		if key == "" {
			return "error: empty key"
		}
		
		if value == "" {
			// Just lookup
			if val, ok := c.KV.Get(key); ok {
				return val.String
			}
			return "key not found: " + key
		}
		
		// Set value
		c.KV.Set(key, value)
		c.KV.Save()
		return "(ok)"
	}
	
	// Handle special commands
	cmd = trimSpace(cmd)
	switch cmd {
	case "keys", "keys *":
		keys := c.KV.Keys()
		if len(keys) == 0 {
			return "(empty)"
		}
		return joinStrings(keys, ", ")
	case "clear":
		return ""
	default:
		// Search in KV
		results := c.KV.Search(cmd)
		if len(results) > 0 {
			return joinStrings(results, ", ")
		}
		return "not found: " + cmd
	}
}

func findIndexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for i := 1; i < len(ss); i++ {
		result += sep + ss[i]
	}
	return result
}

func (c *Core) OpenFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	c.Editor.Load(string(data))
	c.Editor.FilePath = path
	c.CurrentView = ViewEditor
	return nil
}

func (c *Core) SaveFile() error {
	if c.Editor.FilePath == "" {
		return nil
	}
	return os.WriteFile(c.Editor.FilePath, []byte(c.Editor.String()), 0644)
}

func (c *Core) SaveKV() error {
	return c.KV.Save()
}
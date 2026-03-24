package kv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ValueType int

const (
	TypeString ValueType = iota
	TypeNumber
	TypeBool
	TypeArray
	TypeObject
)

type Value struct {
	Type    ValueType
	String  string
	Number  float64
	Bool    bool
	Array   []interface{}
	Object  map[string]interface{}
	Created time.Time
	Updated time.Time
}

type KV struct {
	mu   sync.RWMutex
	data map[string]*Value
	dir  string
}

func New() *KV {
	return &KV{
		data: make(map[string]*Value),
	}
}

func (k *KV) Load(dir string) error {
	k.dir = dir
	file := filepath.Join(dir, "data.kv")
	
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		k.set(key, &Value{
			Type:    TypeString,
			String:  value,
			Created: time.Now(),
			Updated: time.Now(),
		})
	}
	
	return nil
}

func (k *KV) Save() error {
	if k.dir == "" {
		return nil
	}
	
	dir := filepath.Join(k.dir, ".ketsu")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	file := filepath.Join(dir, "data.kv")
	var lines []string
	lines = append(lines, "# Ketsu KV Data")
	lines = append(lines, fmt.Sprintf("# Generated: %s", time.Now().Format(time.RFC3339)))
	lines = append(lines, "")
	
	k.mu.RLock()
	for key, val := range k.data {
		lines = append(lines, fmt.Sprintf("%s=%s", key, val.String))
	}
	k.mu.RUnlock()
	
	return os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
}

func (k *KV) Get(key string) (*Value, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	
	val, ok := k.data[key]
	return val, ok
}

func (k *KV) Set(key string, value string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	k.set(key, &Value{
		Type:    TypeString,
		String:  value,
		Created: time.Now(),
		Updated: time.Now(),
	})
}

func (k *KV) SetNumber(key string, value float64) {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	k.set(key, &Value{
		Type:    TypeNumber,
		Number:  value,
		String:  fmt.Sprintf("%v", value),
		Created: time.Now(),
		Updated: time.Now(),
	})
}

func (k *KV) set(key string, val *Value) {
	if existing, ok := k.data[key]; ok {
		val.Created = existing.Created
	}
	k.data[key] = val
}

func (k *KV) Delete(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.data, key)
}

func (k *KV) Keys() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	
	keys := make([]string, 0, len(k.data))
	for k := range k.data {
		keys = append(keys, k)
	}
	return keys
}

func (k *KV) Search(query string) []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	
	query = strings.ToLower(query)
	var results []string
	
	for key := range k.data {
		if strings.Contains(strings.ToLower(key), query) {
			results = append(results, key)
		}
	}
	
	return results
}

func (k *KV) Increment(key string) float64 {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	val, ok := k.data[key]
	if !ok {
		k.data[key] = &Value{
			Type:    TypeNumber,
			Number:  1,
			String:  "1",
			Created: time.Now(),
			Updated: time.Now(),
		}
		return 1
	}
	
	if val.Type != TypeNumber {
		val.Number = 0
	}
	
	val.Number++
	val.String = fmt.Sprintf("%v", val.Number)
	val.Updated = time.Now()
	
	return val.Number
}

func (k *KV) Decrement(key string) float64 {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	val, ok := k.data[key]
	if !ok {
		k.data[key] = &Value{
			Type:    TypeNumber,
			Number:  -1,
			String:  "-1",
			Created: time.Now(),
			Updated: time.Now(),
		}
		return -1
	}
	
	if val.Type != TypeNumber {
		val.Number = 0
	}
	
	val.Number--
	val.String = fmt.Sprintf("%v", val.Number)
	val.Updated = time.Now()
	
	return val.Number
}
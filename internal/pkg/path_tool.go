package pkg

import (
	"path/filepath"
	"strings"
)

type path struct {
	items []string
}

func NewPath() *path {
	return &path{items: make([]string, 0)}
}

func (p *path) String() string {
	if len(p.items) == 0 {
		return ""
	}
	switch {
	case strings.Contains(p.items[0], ":"):
		return strings.Join(p.items, string(filepath.Separator))
	case p.items[0] == ".":
		return strings.Join(p.items, string(filepath.Separator))
	}
	return string(filepath.Separator) + strings.Join(p.items, string(filepath.Separator))
}

func (p *path) Join(path string) *path {
	if path == "" {
		return p
	}
	path = strings.Trim(path, string(filepath.Separator))
	path = strings.Trim(path, "/")
	path = strings.Trim(path, "\\")
	arr := strings.Split(path, string(filepath.Separator))
	p.items = append(p.items, arr...)
	return p
}

func (p *path) Reset() {
	p.items = p.items[0:0]
}

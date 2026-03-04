// internal/cache/cache.go
package cache

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var ErrExpired = errors.New("cache entry expired")

type Cache struct {
	dir string
	ttl time.Duration
}

func New(dir string, ttl time.Duration) *Cache {
	return &Cache{dir: dir, ttl: ttl}
}

func DefaultDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rtcheck", "data")
}

func (c *Cache) Dir() string {
	return c.dir
}

func (c *Cache) Write(name string, data []byte) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.dir, name), data, 0o644)
}

func (c *Cache) Read(name string) ([]byte, error) {
	path := filepath.Join(c.dir, name)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if time.Since(info.ModTime()) > c.ttl {
		return nil, ErrExpired
	}
	return os.ReadFile(path)
}

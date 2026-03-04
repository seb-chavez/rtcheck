// internal/cache/cache_test.go
package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheDir(t *testing.T) {
	dir := t.TempDir()
	c := New(dir, 24*time.Hour)
	if c.Dir() != dir {
		t.Errorf("Dir() = %q, want %q", c.Dir(), dir)
	}
}

func TestWriteAndReadFresh(t *testing.T) {
	dir := t.TempDir()
	c := New(dir, 24*time.Hour)

	data := []byte(`{"test": true}`)
	if err := c.Write("test.json", data); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := c.Read("test.json")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("Read = %q, want %q", got, data)
	}
}

func TestReadExpired(t *testing.T) {
	dir := t.TempDir()
	c := New(dir, 1*time.Millisecond)

	data := []byte(`{"test": true}`)
	if err := c.Write("test.json", data); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err := c.Read("test.json")
	if err != ErrExpired {
		t.Errorf("Read = %v, want ErrExpired", err)
	}
}

func TestReadMissing(t *testing.T) {
	dir := t.TempDir()
	c := New(dir, 24*time.Hour)

	_, err := c.Read("nonexistent.json")
	if !os.IsNotExist(err) {
		t.Errorf("Read = %v, want os.IsNotExist", err)
	}
}

func TestDefaultDir(t *testing.T) {
	d := DefaultDir()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".rtcheck", "data")
	if d != expected {
		t.Errorf("DefaultDir() = %q, want %q", d, expected)
	}
}

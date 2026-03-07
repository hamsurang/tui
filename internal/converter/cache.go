package converter

import (
	"os"
	"path/filepath"
)

func CacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".tui", "cache")
}

func SavePixelArt(rendered string, name string) (string, error) {
	dir := CacheDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, name+".ansi")
	return path, os.WriteFile(path, []byte(rendered), 0644)
}

func LoadPixelArt(name string) (string, error) {
	path := filepath.Join(CacheDir(), name+".ansi")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

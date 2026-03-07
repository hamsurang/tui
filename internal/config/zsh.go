package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpdateZshrc() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	zshrcPath := filepath.Join(home, ".zshrc")

	// Create .zshrc if it doesn't exist
	if _, err := os.Stat(zshrcPath); os.IsNotExist(err) {
		f, err := os.Create(zshrcPath)
		if err != nil {
			return err
		}
		f.Close()
	}

	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "=== tui-theme startup ===") {
		return nil // Already configured
	}

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	script := fmt.Sprintf(`
# === tui-theme startup ===
export TUI_APPLY=true
if [ "$TUI_APPLY" = "true" ]; then
  %s show
fi
# =========================
`, executable)

	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(script); err != nil {
		return err
	}

	return nil
}

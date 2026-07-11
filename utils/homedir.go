package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveDir resolves the given dir if it uses the telda
func ResolveDir(dir string) (string, error) {
	if strings.HasPrefix(dir, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("couldn't resolve user's HOME directory from %s: %v", dir, err)
		}
		if dir == "~" {
			return homedir, nil
		}
		if strings.HasPrefix(dir, "~/") {
			return filepath.Join(homedir, dir[2:]), nil
		}
		// Handle cases like ~/.ssh-tunnel-manager where the suffix starts with anything other than /
		return filepath.Join(homedir, dir[1:]), nil
	}
	return dir, nil
}

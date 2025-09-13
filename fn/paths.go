package fn

import (
	"os"
	"path/filepath"
	"strings"
)

// baseDir returns the root directory used by Pulse to store its files.
func baseDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".pulse")
	os.MkdirAll(dir, 0o755)
	return dir
}

// receiversDir returns the directory where active receivers are tracked.
func receiversDir() string {
	dir := filepath.Join(baseDir(), "receivers")
	os.MkdirAll(dir, 0o755)
	return dir
}

// sanitize replaces path separators in a name to create safe file names.
func sanitize(name string) string {
	name = strings.ReplaceAll(name, string(filepath.Separator), "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

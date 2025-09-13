package fn

import (
	"fmt"
	"os"
	"path/filepath"
)

// FnStop removes the keep file for the given group to stop its listener.
func FnStop(group string) error {
	file := filepath.Join(receiversDir(), sanitize(group)+".keep")
	if err := os.Remove(file); err != nil {
		return err
	}
	fmt.Println("🛑 Listener stopped for", group)
	return nil
}

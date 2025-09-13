package fn

import (
	"fmt"
	"os"
	"strings"
)

// FnList prints the currently active listening groups.
func FnList() {
	dir := receiversDir()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		fmt.Println("No active listeners")
		return
	}
	fmt.Println("Active groups:")
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".keep") {
			name := strings.TrimSuffix(e.Name(), ".keep")
			fmt.Println(" -", name)
		}
	}
}

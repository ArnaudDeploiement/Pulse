package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"pulse/internal/config"
	"pulse/internal/group"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of active listeners",
	RunE: func(cmd *cobra.Command, args []string) error {
		pidDir := config.PidDir()
		entries, err := os.ReadDir(pidDir)
		if err != nil || len(entries) == 0 {
			fmt.Println()
			fmt.Println(ui.Muted.Render("  No active listeners."))
			fmt.Println()
			return nil
		}

		fmt.Println()
		fmt.Println(ui.Title.Render("  Active Listeners"))

		table := ui.Table{
			Headers: []string{"Group", "PID", "Status"},
			Rows:    [][]string{},
		}

		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".pid") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".pid")

			data, err := os.ReadFile(fmt.Sprintf("%s/%s", pidDir, e.Name()))
			if err != nil {
				continue
			}
			pidStr := strings.TrimSpace(string(data))
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			// Check if process is alive
			process, err := os.FindProcess(pid)
			status := ui.BadgeInactive.Render("dead")
			if err == nil {
				err = process.Signal(syscall.Signal(0))
				if err == nil {
					status = ui.BadgeActive.Render("active")
				}
			}

			// Get member count if group exists
			memberCount := "?"
			if g, err := group.Load(name); err == nil {
				memberCount = fmt.Sprintf("%d", len(g.Members))
			}

			table.Rows = append(table.Rows, []string{
				name,
				pidStr,
				status + "  " + ui.Muted.Render(memberCount+" members"),
			})
		}

		fmt.Println(table.Render())
		return nil
	},
}

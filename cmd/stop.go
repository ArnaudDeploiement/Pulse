package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"pulse/internal/config"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <group>",
	Short: "Stop a listener for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		pidFile := fmt.Sprintf("%s/%s.pid", config.PidDir(), name)

		data, err := os.ReadFile(pidFile)
		if err != nil {
			return fmt.Errorf("no active listener found for group %q", name)
		}

		pidStr := strings.TrimSpace(string(data))
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			os.Remove(pidFile)
			return fmt.Errorf("invalid PID file for group %q", name)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			os.Remove(pidFile)
			return fmt.Errorf("process %d not found", pid)
		}

		if err := process.Signal(os.Interrupt); err != nil {
			os.Remove(pidFile)
			return fmt.Errorf("failed to stop process: %w", err)
		}

		os.Remove(pidFile)
		fmt.Println(ui.Success.Render(fmt.Sprintf("  Listener for %q stopped (PID %d)", name, pid)))
		return nil
	},
}

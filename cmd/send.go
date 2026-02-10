package cmd

import (
	"context"
	"fmt"
	"os"

	"pulse/internal/group"
	"pulse/internal/identity"
	"pulse/internal/transport"
	"pulse/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send <group> <file>",
	Short: "Send a file to all group members",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupName := args[0]
		filePath := args[1]

		// Validate file exists
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("file not found: %s", filePath)
		}
		if info.IsDir() {
			return fmt.Errorf("directories not supported yet, please specify a file")
		}

		// Load group
		g, err := group.Load(groupName)
		if err != nil {
			return err
		}

		if len(g.Members) == 0 {
			fmt.Println(ui.Warning.Render("  No members in group. Add members with: pulse group add " + groupName + " <peerID>"))
			return nil
		}

		// Show confirmation
		fmt.Println()
		fmt.Printf("  %s %s %s %s %s %d %s\n",
			ui.Subtitle.Render("Send"),
			ui.Highlight.Render(info.Name()),
			ui.Muted.Render(fmt.Sprintf("(%s)", formatSize(info.Size()))),
			ui.Muted.Render("to"),
			ui.Subtitle.Render(fmt.Sprintf("%d", len(g.Members))),
			len(g.Members),
			ui.Muted.Render("peer(s) in group"),
		)
		fmt.Println()

		// Load identity
		priv, _, err := identity.LoadPrivateKey()
		if err != nil {
			return fmt.Errorf("loading identity: %w", err)
		}

		// Create progress channel bridging transport events to UI
		progressCh := make(chan transport.SendProgress, len(g.Members))
		uiCh := make(chan ui.PeerResult, len(g.Members))

		// Bridge transport progress to UI
		go func() {
			for p := range progressCh {
				uiCh <- ui.PeerResult{
					PeerID: p.PeerID,
					Ok:     p.Done,
					Err:    p.Err,
				}
			}
			close(uiCh)
		}()

		// Start transfer in background
		ctx := context.Background()
		errCh := make(chan error, 1)
		go func() {
			errCh <- transport.SendFile(ctx, priv, g, filePath, progressCh)
		}()

		// Run progress UI
		model := ui.NewProgress(len(g.Members), uiCh)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			return err
		}

		// Check for transport-level error
		if err := <-errCh; err != nil {
			return err
		}

		return nil
	},
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

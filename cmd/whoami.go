package cmd

import (
	"fmt"

	"pulse/internal/config"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display your PeerID",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		box := ui.InfoBox.Render(
			ui.Subtitle.Render("Your Identity") + "\n\n" +
				ui.KeyValue("PeerID", cfg.PeerID) + "\n" +
				ui.KeyValue("Relay", displayRelay(cfg.DefaultRelay)),
		)
		fmt.Println(box)
		return nil
	},
}

func displayRelay(relay string) string {
	if relay == "" {
		return ui.Muted.Render("(not set)")
	}
	return relay
}

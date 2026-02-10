package cmd

import (
	"fmt"

	"pulse/internal/config"
	"pulse/internal/identity"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Pulse identity and configuration",
	Long:  "Generate a new Ed25519 identity and create the Pulse configuration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.IsInitialized() {
			cfg, err := config.Load()
			if err == nil {
				fmt.Println(ui.Warning.Render("Pulse is already initialized."))
				fmt.Println(ui.KeyValue("PeerID", cfg.PeerID))
				fmt.Println(ui.KeyValue("Config", config.ConfigPath()))
				return nil
			}
		}

		relay, _ := cmd.Flags().GetString("relay")

		result, err := ui.RunSpinner("Generating identity...", func() (string, error) {
			peerID, err := identity.Generate()
			if err != nil {
				return "", err
			}

			cfg := config.Config{
				PeerID:       peerID,
				DefaultRelay: relay,
			}
			if err := config.Save(cfg); err != nil {
				return "", fmt.Errorf("saving config: %w", err)
			}

			s := ui.SuccessBox.Render(
				ui.Success.Render("Pulse initialized!") + "\n\n" +
					ui.KeyValue("PeerID", peerID) + "\n" +
					ui.KeyValue("Config", config.ConfigPath()) + "\n" +
					ui.KeyValue("Key", config.IdentityKeyPath()),
			)
			if relay != "" {
				s += "\n" + ui.KeyValue("Default relay", relay)
			}
			s += "\n\n" + ui.Muted.Render("Share your PeerID with peers so they can add you to groups.")
			return s, nil
		})

		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("relay", "r", "", "Default relay address (multiaddr)")
}

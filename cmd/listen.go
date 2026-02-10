package cmd

import (
	"context"
	"fmt"

	"pulse/internal/group"
	"pulse/internal/identity"
	"pulse/internal/transport"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen <group>",
	Short: "Listen for incoming files from a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupName := args[0]
		storeDir, _ := cmd.Flags().GetString("dir")
		if storeDir == "" {
			storeDir = "./" + groupName
		}

		g, err := group.Load(groupName)
		if err != nil {
			return err
		}

		priv, peerID, err := identity.LoadPrivateKey()
		if err != nil {
			return fmt.Errorf("loading identity: %w", err)
		}

		fmt.Println()
		fmt.Println(ui.KeyValue("PeerID", peerID))
		fmt.Println(ui.KeyValue("Group", groupName))
		fmt.Println(ui.KeyValue("Store", storeDir))
		fmt.Println()

		// Connect and start listening
		var lr *transport.ListenResult
		_, err = ui.RunSpinner("Connecting to relay...", func() (string, error) {
			var listenErr error
			lr, listenErr = transport.Listen(context.Background(), priv, g, storeDir)
			if listenErr != nil {
				return "", listenErr
			}
			return ui.Success.Render("Connected to relay!"), nil
		})
		if err != nil {
			return err
		}

		// Run the interactive listener UI
		return ui.RunListener(groupName, storeDir, lr.Events)
	},
}

func init() {
	listenCmd.Flags().StringP("dir", "d", "", "Directory to store received files (default: ./<group>)")
}

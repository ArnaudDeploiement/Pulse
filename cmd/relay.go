package cmd

import (
	"context"
	"fmt"

	"pulse/internal/transport"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Start a relay server for NAT traversal",
	Long:  "Run a libp2p circuit relay v2 server. Peers use the printed address to connect through NAT/firewalls.",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")

		fmt.Println()
		var info *transport.RelayInfo
		var done <-chan struct{}

		result, err := ui.RunSpinner("Starting relay server...", func() (string, error) {
			var startErr error
			info, done, startErr = transport.StartRelay(context.Background(), port)
			if startErr != nil {
				return "", startErr
			}
			return ui.Success.Render("Relay server running!"), nil
		})
		if err != nil {
			return err
		}
		fmt.Println(result)
		fmt.Println()

		fmt.Println(ui.KeyValue("PeerID", info.PeerID))
		fmt.Println(ui.KeyValue("Port", fmt.Sprintf("%d", port)))
		fmt.Println()

		fmt.Println(ui.Subtitle.Render("  Relay addresses:"))
		for _, addr := range info.Addrs {
			fmt.Printf("  %s\n", ui.Highlight.Render(addr))
		}

		fmt.Println()
		fmt.Println(ui.Muted.Render("  Use one of these addresses with:"))
		fmt.Println(ui.Muted.Render("    pulse init --relay <address>"))
		fmt.Println(ui.Muted.Render("    pulse group create <name> --relay <address>"))
		fmt.Println()
		fmt.Println(ui.Muted.Render("  Press Ctrl+C to stop the relay."))

		<-done
		fmt.Println()
		fmt.Println(ui.Success.Render("  Relay stopped."))
		return nil
	},
}

func init() {
	relayCmd.Flags().IntP("port", "p", 4001, "TCP port to listen on")
}

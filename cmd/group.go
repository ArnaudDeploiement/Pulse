package cmd

import (
	"fmt"

	"pulse/internal/config"
	"pulse/internal/group"
	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var groupCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		relay, _ := cmd.Flags().GetString("relay")

		// Fall back to default relay
		if relay == "" {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			relay = cfg.DefaultRelay
		}
		if relay == "" {
			return fmt.Errorf("relay address required: use --relay or set a default relay with 'pulse init --relay'")
		}

		result, err := ui.RunSpinner(fmt.Sprintf("Creating group %q...", name), func() (string, error) {
			g, err := group.Create(name, relay)
			if err != nil {
				return "", err
			}
			return ui.SuccessBox.Render(
				ui.Success.Render(fmt.Sprintf("Group %q created!", name)) + "\n\n" +
					ui.KeyValue("Protocol", g.Protocol) + "\n" +
					ui.KeyValue("Relay", g.Relay) + "\n\n" +
					ui.Muted.Render("Add members with: pulse group add "+name+" <peerID>"),
			), nil
		})
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var groupAddCmd = &cobra.Command{
	Use:   "add <group> <peerID> [peerID...]",
	Short: "Add member(s) to a group",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		peerIDs := args[1:]

		for _, pid := range peerIDs {
			if err := group.AddMember(name, pid); err != nil {
				fmt.Println(ui.Error.Render(fmt.Sprintf("  Failed to add %s: %s", pid, err)))
				continue
			}
			short := pid
			if len(short) > 20 {
				short = short[:10] + "..." + short[len(short)-8:]
			}
			fmt.Println(ui.Success.Render("  Added ") + ui.Highlight.Render(short) + ui.Muted.Render(" to "+name))
		}
		return nil
	},
}

var groupRemoveCmd = &cobra.Command{
	Use:   "remove <group> <peerID>",
	Short: "Remove a member from a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := group.RemoveMember(args[0], args[1]); err != nil {
			return err
		}
		fmt.Println(ui.Success.Render(fmt.Sprintf("  Removed peer from %q", args[0])))
		return nil
	},
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		groups, err := group.List()
		if err != nil {
			return err
		}

		if len(groups) == 0 {
			fmt.Println(ui.Muted.Render("  No groups. Create one with: pulse group create <name> --relay <addr>"))
			return nil
		}

		fmt.Println()
		fmt.Println(ui.Title.Render("  Groups"))

		table := ui.Table{
			Headers: []string{"Name", "Members", "Relay"},
			Rows:    make([][]string, 0, len(groups)),
		}
		for _, g := range groups {
			relay := g.Relay
			if len(relay) > 40 {
				relay = relay[:37] + "..."
			}
			table.Rows = append(table.Rows, []string{
				g.Name,
				fmt.Sprintf("%d", len(g.Members)),
				relay,
			})
		}
		fmt.Println(table.Render())
		return nil
	},
}

var groupInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show group details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := group.Load(args[0])
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Title.Render("  Group: " + g.Name))
		fmt.Println(ui.KeyValue("Protocol", g.Protocol))
		fmt.Println(ui.KeyValue("Relay", g.Relay))
		fmt.Println(ui.KeyValue("Members", fmt.Sprintf("%d", len(g.Members))))
		fmt.Println()

		if len(g.Members) > 0 {
			fmt.Println(ui.Subtitle.Render("  Members:"))
			for i, m := range g.Members {
				short := m
				if len(short) > 20 {
					short = short[:10] + "..." + short[len(short)-8:]
				}
				fmt.Printf("  %s %s\n", ui.Muted.Render(fmt.Sprintf("%d.", i+1)), ui.Highlight.Render(short))
				fmt.Printf("     %s\n", ui.Muted.Render(m))
			}
		}
		fmt.Println()
		return nil
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := group.Delete(args[0]); err != nil {
			return err
		}
		fmt.Println(ui.Success.Render(fmt.Sprintf("  Group %q deleted.", args[0])))
		return nil
	},
}

func init() {
	groupCreateCmd.Flags().StringP("relay", "r", "", "Relay address (multiaddr)")
	groupCmd.AddCommand(groupCreateCmd, groupAddCmd, groupRemoveCmd, groupListCmd, groupInfoCmd, groupDeleteCmd)
}

package cmds

import (
	"github.com/spf13/cobra"
	"pulse/fn"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Protocols list",
		Run: func(cmd *cobra.Command, args []string) {
			fn.FnList()
		},
	}
	return cmd
}

package cmds

import (
	"fmt"
	"pulse/fn"

	"github.com/spf13/cobra"
)

func IdCmd() *cobra.Command {
	var mode string 

	cmd := &cobra.Command{
		Use:   "id",
		Short: "Get or Add PeerId",
		Run: func(cmd *cobra.Command, args []string) {
		
			switch mode {
			case "id":
				fmt.Println("get id")
			case "add":
				fn.AddPeerId(args[0:])

			default:
				fmt.Println("erreur")
			}
		},
	}

	
	cmd.Flags().StringVarP(&mode, "mode", "m", "id", "Action: id or add")

	return cmd
}

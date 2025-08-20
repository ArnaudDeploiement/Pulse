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
			case "get":
				peerid:=fn.Getid()
				fmt.Printf("Your PeerID : %s", peerid)
			case "add":
				fn.AddPeerId(args[0:])
			default:
				fmt.Println("erreur")
			}
		},
	}

	
	cmd.Flags().StringVarP(&mode, "mode", "m", "get", "Action: get or add")

	return cmd
}

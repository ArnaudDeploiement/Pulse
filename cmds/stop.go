package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"pulse/fn"
)

func StopCmd() *cobra.Command {
	var group string
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop protocol listener",
		Run: func(cmd *cobra.Command, args []string) {
			if group == "" {
				fmt.Printf("Vous devez préciser le nom du groupe pour l'arrêter --name")
				return
			}
			if err := fn.FnStop(group); err != nil {
				fmt.Println("error:", err)
			}
		},
	}
	cmd.Flags().StringVarP(&group, "name", "n", "", "Nom du groupe")
	return cmd
}

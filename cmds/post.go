package cmds

import (
	"fmt"
	"os"
	"pulse/fn"

	"github.com/spf13/cobra"
)

func PostCmd() *cobra.Command{

	var file string
	var protocol string
	var idFile string
	
	cmd := &cobra.Command{
		Use:   "post",
		Short: "Publish data to a protocol",
		Run: func(cmd *cobra.Command, args []string) {
		if protocol=="" || file == "" || idFile ==""  {
				fmt.Printf("You have to specify Protocol, File & idFile path")
				os.Exit(1)
			}
			
			fn.FnPost(protocol,file, idFile)
			
			
		},
	}

	cmd.Flags().StringVarP(&protocol, "protocol", "p", "", "Protocol Path")
	cmd.MarkFlagRequired("protocol")
	cmd.Flags().StringVarP(&file, "file", "f", "", "File Path")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&idFile, "IdFile", "i", "", "IdFile Path")
	cmd.MarkFlagRequired("IdFile")

	return cmd

}
package cmds

import (
	"fmt"
	"os"
	"pulse/fn"

	"github.com/spf13/cobra"
)

func GetCMD() *cobra.Command{
	var storeDir string;
	var protocol string;
	var key string;

		cmd := &cobra.Command{
		Use:   "get",
		Short: "Get data",
		Run: func(cmd *cobra.Command, args []string) {
			if storeDir == "" || protocol == "" || key== ""{
				fmt.Printf("You must specify a repository path --d\nYou must specify the path to a protocol --p\nYou must specify the path to a private key")
				os.Exit(1)
			} 
			fn.FnGet(protocol,storeDir, key)	
		},
	}



	cmd.Flags().StringVarP(&storeDir, "repository", "r", "", "Repository path")
	cmd.Flags().StringVarP(&protocol, "protocol", "p", "", "Protocol path")
	cmd.Flags().StringVarP(&key, "private key", "k", "", "Private key path")
	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("protocol")
	cmd.MarkFlagRequired("private key")

	return cmd

}
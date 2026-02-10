package cmd

import (
	"os"

	"pulse/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Pulse - P2P file sharing",
	Long:  ui.Banner() + ui.Version(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

const rootHelpTmpl = `
{{.Long}}

Usage:
  pulse <command> [flags]

Commands:
{{range .Commands}}{{if .IsAvailableCommand}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}
Flags:
  -h, --help   Show this help

Use "pulse <command> --help" for more information about a command.
`

const subHelpTmpl = `{{if .Long}}{{.Long}}

{{end}}Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(
		initCmd,
		whoamiCmd,
		groupCmd,
		sendCmd,
		listenCmd,
		statusCmd,
		stopCmd,
	)

	rootCmd.SetHelpTemplate(rootHelpTmpl)

	// Set a proper help template for all subcommands
	for _, c := range rootCmd.Commands() {
		c.SetHelpTemplate(subHelpTmpl)
		for _, sc := range c.Commands() {
			sc.SetHelpTemplate(subHelpTmpl)
		}
	}
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

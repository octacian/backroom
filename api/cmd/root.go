package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "backroom",
	Short: "Backroom data aggregator CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

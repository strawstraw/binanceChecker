package main

import (
	"github.com/starwstraw/binanceChecker/client"
	"github.com/starwstraw/binanceChecker/server"

	"github.com/spf13/cobra"
)

func main() {
	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}

	var pair string
	var cmdRater = &cobra.Command{
		Use:   "rate",
		Short: "get rate",
		Run: func(cmd *cobra.Command, args []string) {
			client.GetCurrentRate(pair)
		},
	}
	cmdRater.Flags().StringVarP(&pair, "pair", "p", "", "set token pair for lookup")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdServer, cmdRater)
	rootCmd.Execute()
}

package main

import (
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kvdb",
		Short: "A key-value database",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the KVDB server.",
		Run: func(cmd *cobra.Command, args []string) {
			err := StartServer()
			if err != nil {
				println("error: " + err.Error())
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(serveCmd)

	err := rootCmd.Execute()
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

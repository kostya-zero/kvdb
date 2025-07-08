package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var port int
	var file string

	rootCmd := &cobra.Command{
		Use:   "kvdb",
		Short: "A key-value database",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the KVDB server.",
		Run: func(cmd *cobra.Command, args []string) {
			err := StartServer(port, file)
			if err != nil {
				println("error: " + err.Error())
				os.Exit(1)
			}
		},
	}

	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "The port number to use.")
	serveCmd.Flags().StringVarP(&file, "file", "f", "", "The file to use for database.")

	rootCmd.AddCommand(serveCmd)

	err := rootCmd.Execute()
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

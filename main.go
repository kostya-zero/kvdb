package main

import (
	"os"

	"github.com/spf13/cobra"
)

func GetEnv(name string, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" {
		return defaultValue
	}
	return variable
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "kvdb",
		Short: "A key-value database",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the KVDB server.",
		Run:   serveCommand,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of KVDB",
		Run:   versionCommand,
	}

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	err := rootCmd.Execute()
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

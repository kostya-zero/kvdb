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

	overviewCmd := &cobra.Command{
		Use:   "overview",
		Short: "Inspect the database with tree view.",
		Run:   overviewCommand,
	}

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(overviewCmd)

	err := rootCmd.Execute()
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

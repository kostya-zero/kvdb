package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	_ "net/http/pprof"
)

func GetEnv(name string, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" {
		return defaultValue
	}
	return variable
}

func main() {
	go func() {
		println(http.ListenAndServe("localhost:6060", nil))
	}()

	rootCmd := &cobra.Command{
		Use:   "kvdb",
		Short: "A key-value database",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the KVDB server.",
		Run: func(cmd *cobra.Command, args []string) {
			port := GetEnv("KVDB_PORT", "5511")
			portInt, err := strconv.Atoi(port)
			if err != nil {
				LogWarn("Invalid port in KVDB_PORT. Falling back to 5511.")
				portInt = 5511
			}

			file := GetEnv("KVDB_DATABASE", "")

			interval := GetEnv("KVDB_SAVE_INTERVAL", "60000")
			intervalInt, err := strconv.Atoi(interval)
			if err != nil {
				LogWarn("Invalid interval in KVDB_SAVE_INTERVAL. Falling back to 60000.")
				intervalInt = 60000
			}

			err = StartServer(portInt, file, intervalInt)
			if err != nil {
				println("error: " + err.Error())
				os.Exit(1)
			}
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of KVDB",
		Run: func(cmd *cobra.Command, args []string) {
			println("KVDB version: " + version)
		},
	}

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	err := rootCmd.Execute()
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

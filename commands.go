package main

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func serveCommand(cmd *cobra.Command, args []string) {
	port := GetEnv("KVDB_PORT", "5511")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		LogWarn("cli", "Invalid port in KVDB_PORT. Falling back to 5511.")
		portInt = 5511
	}

	file := GetEnv("KVDB_DATABASE", "")

	interval := GetEnv("KVDB_SAVE_INTERVAL", "60000")
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		LogWarn("cli", "Invalid interval in KVDB_SAVE_INTERVAL. Falling back to 60000.")
		intervalInt = 60000
	}

	err = StartServer(portInt, file, intervalInt)
	if err != nil {
		println("error: " + err.Error())
		os.Exit(1)
	}
}

func versionCommand(cmd *cobra.Command, args []string) {
	println("KVDB version: " + version)
}

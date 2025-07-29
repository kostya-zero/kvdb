package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
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

func overviewCommand(cmd *cobra.Command, args []string) {
	file := GetEnv("KVDB_DATABASE", "")

	if file == "" {
		println("Specify path to the database using KVDB_DATABASE environment variable.")
		os.Exit(1)
	}

	db := NewDatabase(file)
	err := db.LoadFromFile()
	if err != nil {
		println("An error occured while loading database: " + err.Error())
		os.Exit(1)
	}

	if len(db.Maps) == 0 {
		println("Database is empty.")
		return
	}

	filename := filepath.Base(file)

	tree := treeprint.NewWithRoot(filename)

	for dbMap, keys := range db.Maps {
		newBranch := tree.AddBranch(dbMap)

		if len(keys) == 0 {
			newBranch.AddNode("... empty map ...")
			continue
		}

		for key, value := range keys {
			newBranch.AddNode(fmt.Sprintf("%s: %s", key, value))
		}
	}

	fmt.Println(tree.String())
}

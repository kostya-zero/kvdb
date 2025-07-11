package main

import (
	"fmt"
	"time"
)

func LogInfo(msg string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] \x1b[94mINFO\x1b[0m: " + msg)
}

func LogError(msg string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] \x1b[91mERROR\x1b[0m: " + msg)
}

func LogWarn(msg string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] \x1b[93mWARN\x1b[0m: " + msg)
}

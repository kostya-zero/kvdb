package main

import (
	"fmt"
	"time"
)

func LogInfo(scope string, msg string) {
	fmt.Printf("["+time.Now().Format("2006-01-02 15:04:05")+"] \x1b[94mINFO\x1b[0m(%s): %s\n", scope, msg)
}

func LogError(scope string, msg string) {
	fmt.Printf("["+time.Now().Format("2006-01-02 15:04:05")+"] \x1b[91mERROR\x1b[0m(%s): %s\n", scope, msg)
}

func LogWarn(scope string, msg string) {
	fmt.Printf("["+time.Now().Format("2006-01-02 15:04:05")+"] \x1b[93mWARN\x1b[0m(%s): %s\n", scope, msg)
}

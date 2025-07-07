package main

import (
	"fmt"
	"time"
)

type Logger struct{}

func LogInfo(msg string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] \x1b[94mINFO\x1b[0m: " + msg)
}

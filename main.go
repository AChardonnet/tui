package main

import (
	"log/slog"
)

const ()

func main() {
	logChannel := make(chan logMsg, 100)
	logFile, err := setupLogger(logChannel)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	slog.Debug("Starting ...")

	tui(logChannel)
}

package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/octacian/backroom/api/cmd"
	"github.com/octacian/backroom/api/config"
	slogmulti "github.com/samber/slog-multi"
)

func main() {
	// Configure multi logger
	stderr := os.Stderr
	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(logfile, &slog.HandlerOptions{}),
			tint.NewHandler(stderr, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.Kitchen,
			}),
		),
	)

	slog.SetDefault(logger)

	// Initialize configuration
	config.Init()

	// Initialize command line interface
	cmd.Execute()
}

package main

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/file-share/app/api/handlers"
	"github.com/FlameInTheDark/file-share/foundation/logs"
)

func ApiCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "start api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log",
				Value:   "prod",
				Usage:   "set log level (prod, debug)",
				EnvVars: []string{"LOG_LEVEL"},
			},
		},
		Action: func(c *cli.Context) error {
			// Create logger
			var logger *zap.Logger
			switch c.String("log") {
			case "debug":
				logger = logs.CreateLoggerDebug()
			case "prod":
				logger = logs.CreateLogger()
			default:
				logger = logs.CreateLogger()
			}
			defer logger.Sync()

			logger.Info("service is starting", zap.String("app-name", c.App.Name), zap.String("app-version", c.App.Version))
			return handlers.Run(logger)
		},
	}
}

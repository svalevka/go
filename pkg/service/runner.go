package service

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"gitlab.com/laserdigital/platform/go/pkg/config"
	"gitlab.com/laserdigital/platform/go/pkg/tasks"
)

// Runner contains the runtime environment of the Service to be configured.
type Runner struct {
	// Tasks contains the concurrent Tasks to be run for the Service, at least
	// one must be defined.
	Tasks *tasks.Runner

	// Logger is a structures logger configured with service metadata.
	Logger *slog.Logger
}

// Run is a generic function that automatically reads a JSON or YAML
// configuration file, then invokes a setup function for the Service, lastly
// running the configured tasks.
func Run[CONFIG any](serviceName string, setup func(context.Context, *Runner, *CONFIG) error) int {
	ctx := context.Background()

	configFile := flag.String("config", "config.yml", "path to JSON or YAML configuration file")
	flag.Parse()

	cfg := new(CONFIG)
	err := config.FromFile(*configFile, cfg)
	if err != nil {
		return exitError(2, "Config: %s", err)
	}

	log := getLogger(cfg).With(slog.String("service", serviceName))

	rn := &Runner{
		Tasks: &tasks.Runner{
			TaskStarting: func(ts *tasks.TasksStatus) {
				log.Info("task starting...", slog.String("task", ts.Name))
			},
			TaskStopped: func(ts *tasks.TasksStatus) {
				log.Info("task stopped", slog.String("task", ts.Name))
			},
			TaskFailed: func(ts *tasks.TasksStatus, err error) {
				log.Error("task failed", slog.String("task", ts.Name), slog.String("error", err.Error()))
			},
		},
		Logger: log,
	}

	err = setup(ctx, rn, cfg)
	if err != nil {
		return exitError(1, "Setup: %s", err)
	}

	err = rn.Tasks.Run(ctx)
	if err != nil {
		return exitError(1, "Run: %s", err)
	}

	return 0
}

// getLogger returns a structured logger, either configured by the service
// configuration using the GetLogger interface, or a default of JSON/DEBUG.
func getLogger(cfg any) *slog.Logger {
	if gl, ok := cfg.(interface {
		GetLogger() *slog.Logger
	}); ok {
		return gl.GetLogger()
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func exitError(code int, format string, args ...any) int {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	return code
}

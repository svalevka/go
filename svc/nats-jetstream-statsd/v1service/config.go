package v1service

import (
	"context"
	"log/slog"

	"gitlab.com/laserdigital/platform/go/pkg/config/common"
	"gitlab.com/laserdigital/platform/go/pkg/service"
)

type Config struct {
	common.Logs `yaml:"logs"`

	Streams []*Stream `yaml:"streams"`

	StatsD *StatsD `yaml:"statsd"`
}

type Stream struct {
	Name string `yaml:"name"`

	NATS common.NATS `yaml:"nats"`
}

type StatsD struct {
	Host string `yaml:"host"`
}

func New(ctx context.Context, svc *service.Runner, cfg *Config) error {
	events := make(chan *Event, 1024)

	for _, stream := range cfg.Streams {
		svc.Tasks.Add(&Events{
			Stream: stream,
			Target: events,
			Logger: svc.Logger.With(slog.String("task", "Events"), slog.String("stream", stream.Name)),
		})
	}

	svc.Tasks.Add(&Stats{
		StatsD: cfg.StatsD,
		Source: events,
		Logger: svc.Logger.With(slog.String("task", "Stats")),
	})

	return nil
}

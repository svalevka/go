package v1service

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"

	"gitlab.com/laserdigital/platform/go/pkg/config/common"
	"gitlab.com/laserdigital/platform/go/pkg/encoding"
	"gitlab.com/laserdigital/platform/go/pkg/net/http/api"
	"gitlab.com/laserdigital/platform/go/pkg/service"
	"gitlab.com/laserdigital/platform/go/pkg/tasks"
)

// Config contains the configuration used to configure the systemd-service-ui
// web application and server.
type Config struct {
	common.Logs `yaml:"logs"`

	// Listen configures the host:port where the app will accept new
	// connections.
	Listen []string `yaml:"listen"`

	// Services is an array of regular expression controlling what services
	// the systemd-service-ui is allowed to manage, no other services can
	// be managed without this.
	Services []string `yaml:"services"`
}

// New initializes the service runner for the system.
func New(ctx context.Context, svc *service.Runner, cfg *Config) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("could not get hostname: %w", err)
	}

	managed := []*regexp.Regexp{}

	for _, expr := range cfg.Services {
		re, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("invalid regular expression %q: %w", expr, err)
		}

		managed = append(managed, re)
	}

	systemd, err := NewDbus(ctx, managed)
	if err != nil {
		return fmt.Errorf("could not connect to systemd: %w", err)
	}

	app := &App{
		Hostname: hostname,
		Systemd:  systemd,
		Logger:   svc.Logger,
	}

	rest := &API{
		Hostname: hostname,
		Systemd:  systemd,
		Logger:   svc.Logger,
	}

	r := chi.NewRouter()
	ra := api.From(&encoding.JSON{}, svc.Logger, r)

	r.Route("/", app.Routes)
	ra.Route("/api", rest.Routes)

	for _, addr := range cfg.Listen {
		svc.Tasks.Add(&tasks.HTTPServer{
			Name:    "App",
			Addr:    addr,
			Handler: r,
		})
	}

	return nil
}

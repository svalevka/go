package v1service

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.com/laserdigital/platform/go/pkg/net/http/web/assets"
	"gitlab.com/laserdigital/platform/go/svc/systemd-service-ui/v1service/views"
)

// App implements the Web Application portion of systemd-service-ui.
type App struct {
	// Hostname is the name of the host running systemd-service-ui, as
	// displayed on the index page of the web app.
	Hostname string

	// Systemd is the connection implementation to the host Dbus.
	Systemd Systemd

	// Logger is the optional logger where debugging and errors are written to.
	Logger *slog.Logger
}

// Routes attaches the routes of the Web Application to a Chi Router.
func (a *App) Routes(r chi.Router) {
	r.NotFound(a.handle(a.NotFound))
	r.MethodNotAllowed(a.handle(a.MethodNotAllowed))

	r.Get("/", a.handle(a.ListServices))
	r.Post("/services", a.handle(a.ServiceAction))

	r.Handle("/assets/common/js/*", http.StripPrefix("/assets/common/js", http.FileServer(http.FS(assets.JS()))))
}

func (a *App) NotFound(w http.ResponseWriter, r *http.Request) (views.View, error) {
	return &views.Error{
		Status:  http.StatusNotFound,
		Message: "Page or Resource not found.",
	}, nil
}

func (a *App) MethodNotAllowed(w http.ResponseWriter, r *http.Request) (views.View, error) {
	return &views.Error{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method Not Allowed for Page or Resource.",
	}, nil
}

func (a *App) ListServices(w http.ResponseWriter, r *http.Request) (views.View, error) {
	services, err := a.Systemd.ListServices(r.Context())
	if err != nil {
		return nil, err
	}

	return &views.ListServices{
		Hostname: a.Hostname,
		Services: services,
	}, nil
}

func (a *App) ServiceAction(w http.ResponseWriter, r *http.Request) (views.View, error) {
	ctx := r.Context()

	service := r.PostFormValue("service")

	switch r.PostFormValue("action") {
	case "start":
		err := a.Systemd.StartService(ctx, service)
		if err != nil {
			return nil, err
		}

	case "stop":
		err := a.Systemd.StopService(ctx, service)
		if err != nil {
			return nil, err
		}

	case "restart":
		err := a.Systemd.RestartService(ctx, service)
		if err != nil {
			return nil, err
		}

	default:
		return &views.InlineError{
			Message: "Unknown service action",
		}, nil
	}

	// fetch service again with new state to render component again,
	svc, err := a.Systemd.GetService(ctx, service)
	if err != nil {
		return nil, err
	}

	return &views.ServiceControl{Service: svc}, nil
}

func (a *App) handle(fn func(w http.ResponseWriter, r *http.Request) (views.View, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := fn(w, r)
		if err != nil {
			a.Logger.Error("an unexpected error occurred", slog.String("error", err.Error()))

			view = &views.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		if view != nil {
			err = views.Render(w, r, view)
			if err != nil {
				a.Logger.Error("could not render template", slog.String("error", err.Error()))
			}
		}
	}
}

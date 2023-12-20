package v1service

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/svalevka/go/pkg/net/http/api"
	v1 "github.com/svalevka/go/svc/systemd-service-ui/v1"
)

// API is the JSON API equivalent of the Web Application.
type API struct {
	// Hostname is the name of the host running systemd-service-ui, as
	// displayed on the list services api call.
	Hostname string

	Systemd Systemd

	Logger *slog.Logger
}

// Routes attaches the routes of the Web Application to a Chi Router.
func (a *API) Routes(app *api.App) {
	app.ErrorHandler = func(err error) any {
		body := &v1.ErrorWrapper{}

		switch err := err.(type) {
		case *v1.Error:
			body.Error = err

		default:
			a.Logger.Error("an unexpected error occurred", slog.String("error", err.Error()))

			body.Error = &v1.Error{
				Code:    "Unknown",
				Message: "An unknown error occurred, check the logs.",
			}
		}

		return body
	}

	app.NotFound(&v1.ErrorWrapper{
		Error: &v1.Error{Code: "NotFound", Message: "Resource not found."},
	})

	app.MethodNotAllowed(&v1.ErrorWrapper{
		Error: &v1.Error{Code: "MethodNotAllowed", Message: "Method Not Allowed for Resource."},
	})

	api.Get(app, "/services", a.ListServices)
	api.Post(app, "/services/{service}:start", a.StartService)
	api.Post(app, "/services/{service}:restart", a.RestartService)
	api.Post(app, "/services/{service}:stop", a.StopService)
}

func (a *API) ListServices(ctx context.Context, req *api.Request[api.None]) (*api.Response[v1.ListServicesRes], error) {
	services, err := a.Systemd.ListServices(ctx)
	if err != nil {
		return nil, err
	}

	return &api.Response[v1.ListServicesRes]{
		Body: &v1.ListServicesRes{
			Hostname: a.Hostname,
			Services: services,
		},
	}, nil
}

func (a *API) StartService(ctx context.Context, req *api.Request[api.None]) (*api.Response[api.None], error) {
	service, _ := url.QueryUnescape(api.URLParam(ctx, "service"))

	err := a.Systemd.StartService(ctx, service)
	if err != nil {
		return nil, err
	}

	return &api.Response[api.None]{
		StatusCode: http.StatusNoContent,
	}, nil
}

func (a *API) StopService(ctx context.Context, req *api.Request[api.None]) (*api.Response[api.None], error) {
	service, _ := url.QueryUnescape(api.URLParam(ctx, "service"))

	err := a.Systemd.RestartService(ctx, service)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *API) RestartService(ctx context.Context, req *api.Request[api.None]) (*api.Response[api.None], error) {
	service, _ := url.QueryUnescape(api.URLParam(ctx, "service"))

	err := a.Systemd.StopService(ctx, service)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

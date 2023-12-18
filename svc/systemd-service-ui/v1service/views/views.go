package views

import (
	"net/http"

	v1 "gitlab.com/laserdigital/platform/go/svc/systemd-service-ui/v1"
)

// View is implemented by all views exposed by this package.
type View interface {
	TemplateName() string
}

type ListServices struct {
	Hostname string
	Services []*v1.Service
}

func (l *ListServices) TemplateName() string {
	return "list_services.html"
}

type ServiceControl struct {
	*v1.Service
}

func (s *ServiceControl) TemplateName() string {
	return "service_control.html"
}

type Error struct {
	Status  int
	Message string
}

func (e *Error) StatusCode() int {
	if e.Status > 0 {
		return e.Status
	}

	return http.StatusInternalServerError
}

func (e *Error) TemplateName() string {
	return "error.html"
}

type InlineError struct {
	Message string
}

func (i *InlineError) TemplateName() string {
	return "inline_error.html"
}

package v1service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"

	v1 "gitlab.com/laserdigital/platform/go/svc/systemd-service-ui/v1"
)

// Systemd abstracts operations against the host Systemd via Dbus.
type Systemd interface {
	ListServices(ctx context.Context) (v1.Services, error)

	GetService(ctx context.Context, service string) (*v1.Service, error)

	StartService(ctx context.Context, service string) error

	RestartService(ctx context.Context, service string) error

	StopService(ctx context.Context, service string) error
}

// Dbus is a Systemd implementation that is backed directly by Dbus.
type Dbus struct {
	managed []*regexp.Regexp
	conn    *dbus.Conn
}

func NewDbus(ctx context.Context, managed []*regexp.Regexp) (*Dbus, error) {
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &Dbus{conn: conn, managed: managed}, nil
}

func (d *Dbus) Close() error {
	d.conn.Close()
	return nil
}

func (d *Dbus) ListServices(ctx context.Context) (v1.Services, error) {
	units, err := d.conn.ListUnitsContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list systemd units: %w", err)
	}

	svc := v1.Services{}

	for _, unit := range units {
		if d.isManagedService(unit.Name) {
			svc = append(svc, &v1.Service{
				Name:        unit.Name,
				Description: unit.Description,
				Running:     unit.ActiveState == "active" && unit.SubState == "running",
			})
		}
	}

	sort.Sort(svc)

	return svc, nil
}

func (d *Dbus) GetService(ctx context.Context, service string) (*v1.Service, error) {
	svcs, err := d.ListServices(ctx)
	if err != nil {
		return nil, err
	}

	for _, svc := range svcs {
		if svc.Name == service {
			return svc, nil
		}
	}

	return nil, fmt.Errorf("service %q not found", service)
}

func (d *Dbus) StartService(ctx context.Context, service string) error {
	if !d.isManagedService(service) {
		return fmt.Errorf("service %q is not found", service)
	}

	reply := make(chan string)
	_, err := d.conn.StartUnitContext(ctx, service, "fail", reply)
	if err != nil {
		return err
	}

	status := <-reply
	if status != "done" {
		return fmt.Errorf("expected status done, got %q", status)
	}

	return nil
}

func (d *Dbus) RestartService(ctx context.Context, service string) error {
	if !d.isManagedService(service) {
		return fmt.Errorf("service %q is not found", service)
	}

	reply := make(chan string)
	_, err := d.conn.RestartUnitContext(ctx, service, "fail", reply)
	if err != nil {
		return err
	}

	status := <-reply
	if status != "done" {
		return fmt.Errorf("expected status done, got %q", status)
	}

	return nil
}

func (d *Dbus) StopService(ctx context.Context, service string) error {
	if !d.isManagedService(service) {
		return fmt.Errorf("service %q is not found", service)
	}

	reply := make(chan string)
	_, err := d.conn.StopUnitContext(ctx, service, "fail", reply)
	if err != nil {
		return err
	}

	status := <-reply
	if status != "done" {
		return fmt.Errorf("expected status done, got %q", status)
	}

	return nil
}

func (d *Dbus) isManagedService(v string) bool {
	if !strings.HasSuffix(v, ".service") {
		return false
	}

	for _, re := range d.managed {
		if re.MatchString(v) {
			return true
		}
	}

	return false
}

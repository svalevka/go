package main

import (
	"os"

	"gitlab.com/laserdigital/platform/go/pkg/service"
	"gitlab.com/laserdigital/platform/go/svc/systemd-service-ui/v1service"
)

func main() {
	os.Exit(service.Run("systemd-service-ui.v1", v1service.New))
}

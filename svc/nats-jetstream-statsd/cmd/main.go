package main

import (
	"os"

	"gitlab.com/laserdigital/platform/go/pkg/service"
	"gitlab.com/laserdigital/platform/go/svc/nats-jetstream-statsd/v1service"
)

func main() {
	os.Exit(service.Run("nats-jetstream-statsd.v1", v1service.New))
}

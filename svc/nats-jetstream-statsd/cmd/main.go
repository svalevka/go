package main

import (
	"os"

	"github.com/svalevka/go/pkg/service"
	"github.com/svalevka/go/svc/nats-jetstream-statsd/v1service"
)

func main() {
	os.Exit(service.Run("nats-jetstream-statsd.v1", v1service.New))
}

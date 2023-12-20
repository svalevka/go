package main

import (
	"os"

	"github.com/svalevka/go/pkg/service"
	"github.com/svalevka/go/svc/systemd-service-ui/v1service"
)

func main() {
	os.Exit(service.Run("systemd-service-ui.v1", v1service.New))
}

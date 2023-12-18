package common

import (
	"github.com/nats-io/nats.go"
)

// NATS contains configuration common to all NATS clients.
type NATS struct {
	// Servers is an array of at least one NATS server to connect to, depending
	// on client implementation, more servers will be discovered.
	Servers []string `yaml:"servers"`

	// Username optionally configures the username to authenticate as, no
	// authentication will take place if empty.
	Username string `yaml:"username"`

	// Password optionally configures the password to authenticate with, no
	// authentication will take place if empty.
	Password string `yaml:"password"`
}

// Connect returns a NATS client configured by NATS, where clientName is used
// to identify this connection to the NATS server for debugging.
func (n *NATS) Connect(clientName string) (*nats.Conn, error) {
	conn, err := nats.Options{
		AllowReconnect: true,
		Servers:        n.Servers,
		User:           n.Username,
		Password:       n.Password,
		Name:           clientName,

		MaxReconnect:       nats.DefaultMaxReconnect,
		ReconnectWait:      nats.DefaultReconnectWait,
		ReconnectJitter:    nats.DefaultReconnectJitter,
		ReconnectJitterTLS: nats.DefaultReconnectJitterTLS,
		Timeout:            nats.DefaultTimeout,
		PingInterval:       nats.DefaultPingInterval,
		MaxPingsOut:        nats.DefaultMaxPingOut,
		SubChanLen:         nats.DefaultMaxChanLen,
		ReconnectBufSize:   nats.DefaultReconnectBufSize,
		DrainTimeout:       nats.DefaultDrainTimeout,
		FlusherTimeout:     nats.DefaultFlusherTimeout,
	}.Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

package v1service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/jsm.go/api/jetstream/advisory"
	server "github.com/nats-io/jsm.go/api/server/advisory"
	"github.com/nats-io/nats.go"
)

// Event contains the metadata unmarshaled from the JetStream
// ConsumerInfoResponse API advisory received that contains information about
// a consumer and their position in the stream.
type Event struct {
	// Client contains information about the account and connection to NATS.
	Client *server.ClientInfoV1

	// Consumer contains information about the stream and the consumers
	// sequence.
	Consumer *api.JSApiConsumerInfoResponse
}

// Events is a task that is invoked for every stream being monitored to
// receive JetStream Advisory events and write them to a channel.
type Events struct {
	Stream *Stream

	Target chan *Event

	// Logger optionally configures the logger where debug information is
	// written to.
	Logger *slog.Logger
}

func (e *Events) TaskName() string {
	return "Events(" + e.Stream.Name + ", " + e.Stream.NATS.Username + ")"
}

func (e *Events) RunTask(ctx context.Context) error {
	conn, err := e.Stream.NATS.Connect("nats-jetstream-statsd.v1")
	if err != nil {
		return fmt.Errorf("nats: %w", err)
	}
	defer conn.Close()

	msgs := make(chan *nats.Msg, 1024)

	sub, err := conn.ChanSubscribe("$JS.EVENT.ADVISORY.>", msgs)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case msg := <-msgs:
			e.handleMsg(msg)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *Events) handleMsg(msg *nats.Msg) {
	kind, body, _ := api.ParseMessage(msg.Data)

	e.Logger.Debug("jetstream event advisory received", slog.String("kind", kind))

	event, ok := body.(*advisory.JetStreamAPIAuditV1)
	if !ok {
		// ignore other events we may receive, it should only ever be
		// this one.
		return
	}

	kind, body, _ = api.ParseMessage([]byte(event.Response))

	e.Logger.Debug("jetstream api audit received", slog.String("kind", kind))

	switch body := body.(type) {
	case *api.JSApiConsumerInfoResponse:
		e.Target <- &Event{
			Client:   &event.Client,
			Consumer: body,
		}
	}
}

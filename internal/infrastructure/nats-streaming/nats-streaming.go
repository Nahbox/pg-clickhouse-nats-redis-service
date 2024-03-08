package nats_streaming

import (
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

func New() (stan.Conn, error) {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return nil, err
	}

	log.Info("nats-streaming connection established")

	return sc, nil
}

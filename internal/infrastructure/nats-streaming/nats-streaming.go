package nats_streaming

import (
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

const (
	clusterID = "test-cluster"
)

func New(clientID string, natsURL string) (stan.Conn, error) {
	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(natsURL),
	)
	if err != nil {
		return nil, err
	}

	log.Info("nats-streaming connection established")

	return sc, nil
}

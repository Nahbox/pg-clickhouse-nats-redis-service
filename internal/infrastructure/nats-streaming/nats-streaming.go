package nats_streaming

import "github.com/nats-io/stan.go"

func New() (stan.Conn, error) {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return nil, err
	}

	return sc, nil
}

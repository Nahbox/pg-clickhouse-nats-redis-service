package nats_streaming

import (
	"encoding/json"

	"github.com/nats-io/stan.go"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/msgbroker"
)

type Repo struct {
	sc stan.Conn
}

func NewMsgbRepo(sc stan.Conn) msgbroker.Repository {
	return &Repo{sc}
}

func (r *Repo) Publish(data *msgbroker.Log) error {
	subject := "test"

	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := r.sc.Publish(subject, msg); err != nil {
		//log.Fatalf("error publishing: %v", err)
		return err
	}

	//log.Printf("published from: %s", v.Name())
	return nil
}

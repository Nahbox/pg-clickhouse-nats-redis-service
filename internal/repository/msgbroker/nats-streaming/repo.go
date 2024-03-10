package nats_streaming

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/logs"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/msgbroker"
)

type Repo struct {
	sc stan.Conn
}

func NewRepo(sc stan.Conn) msgbroker.Repository {
	return &Repo{sc}
}

func (r *Repo) Publish(data *logs.Log) error {
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

// ReadAsync reads in the background each new message from subject in out channel
// ReadAsync is the non-blocking function
func (r *Repo) ReadAsync(ctx context.Context, subjectID string, out chan<- *logs.Log) error {
	sub, err := r.sc.Subscribe(subjectID, func(m *stan.Msg) {
		var msg logs.Log
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			return
		}
		out <- &msg
		fmt.Println("aaaaa")
	}, stan.DeliverAllAvailable())

	if err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			r.sc.Close()
		}
	}()

	return nil
}

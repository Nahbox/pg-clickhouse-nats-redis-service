package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/logs"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/msgbroker"
)

type Service struct {
	LogsRepo      logs.Repository
	MsgBrokerRepo msgbroker.Repository
}

func NewService(logsRepo logs.Repository, msgBrokerRepo msgbroker.Repository) *Service {
	return &Service{
		LogsRepo:      logsRepo,
		MsgBrokerRepo: msgBrokerRepo,
	}
}

func batchLogs(ctx context.Context, logsCh <-chan *logs.Log, maxItems int, maxTimeout time.Duration) chan []logs.Log {
	batches := make(chan []logs.Log)

	sendBatch := func(batch []logs.Log) {
		if len(batch) > 0 {
			batches <- batch
		}
	}

	go func() {
		defer func() {
			close(batches)
			fmt.Println("goroutine finished")

		}()

		for keepGoing := true; keepGoing; {
			var batch []logs.Log

			expire := time.After(maxTimeout)

			brflag := false

			for {
				if brflag {
					break
				}
				select {
				case <-ctx.Done():
					fmt.Println("ctx.Done()")
					keepGoing = false
					sendBatch(batch)
					return

				case value, ok := <-logsCh:
					if !ok {
						fmt.Println("closed channel")

						keepGoing = false
						sendBatch(batch)
						return
					}

					fmt.Println("id = ", value.Id)

					batch = append(batch, *value)
					if len(batch) == maxItems {
						sendBatch(batch)
						brflag = true
						break
					}

				case <-expire:
					sendBatch(batch)
					brflag = true
					break
				}
			}
		}
	}()

	return batches
}

func (s *Service) SaveLogs(ctx context.Context) error {
	logsCh := make(chan *logs.Log)

	err := s.MsgBrokerRepo.ReadAsync(ctx, "test", logsCh)
	if err != nil {
		return err
	}

	batches := batchLogs(ctx, logsCh, 2, 3*time.Second)
	for batch := range batches {
		err := s.LogsRepo.AddBatch(ctx, batch)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("SaveLogs finished")

	return nil
}

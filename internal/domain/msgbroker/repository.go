package msgbroker

type Repository interface {
	Publish(data *Log) error
}

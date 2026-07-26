package queue

import (
	"github.com/nats-io/nats.go"
)

func Connect(url string) (*nats.Conn, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

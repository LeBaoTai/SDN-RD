package nats

import (
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn *nats.Conn
}

func NewSubscriber(url string) (*Subscriber, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		conn: conn,
	}, nil
}

func (s *Subscriber) SubscribeQueue(subject, queueGroup string, handler func(data []byte)) error {
	_, err := s.conn.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}

func (s *Subscriber) Close() {
	s.conn.Close()
}

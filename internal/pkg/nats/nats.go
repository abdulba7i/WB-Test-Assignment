package nats

import (
	"fmt"
	"l0/internal/config"

	"github.com/nats-io/stan.go"
)

type Nats struct {
	connection stan.Conn
}

func New(cfg config.NatsStreaming, clientID string) (*Nats, error) {
	natsURL := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)
	sc, err := stan.Connect(cfg.ClusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return &Nats{connection: sc}, nil
}

func (n *Nats) Publish(topic string, data []byte) error {
	return n.connection.Publish(topic, data)
}

func (n *Nats) Consume(topic string, handler func(msg *stan.Msg), opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return n.connection.Subscribe(topic, handler, opts...)
}

func (n *Nats) Close() {
	n.connection.Close()
}

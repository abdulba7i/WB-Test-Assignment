package nats

import (
	"fmt"
	"l0/internal/config"

	"github.com/nats-io/stan.go"
)

type Nats struct {
	connection stan.Conn
}

// New создает новое подключение к NATS Streaming
func New(cfg config.NatsStreaming, clientID string) (*Nats, error) {
	natsURL := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)
	sc, err := stan.Connect(cfg.ClusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return &Nats{connection: sc}, nil
}

// Publish отправляет сообщение в указанный канал
func (n *Nats) Publish(topic string, data []byte) error {
	return n.connection.Publish(topic, data)
}

// Consume подписывается на канал с обработчиком сообщений
func (n *Nats) Consume(topic string, handler func(msg *stan.Msg), opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return n.connection.Subscribe(topic, handler, opts...)
}

// Close закрывает соединение с NATS Streaming
func (n *Nats) Close() {
	n.connection.Close()
}

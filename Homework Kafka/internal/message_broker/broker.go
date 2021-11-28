package message_broker

import "context"

type (
	MessageBroker interface {
		Connect(ctx context.Context) error
		Close() error

		Cache() CacheBroker // Сюда отправляются команды, связанные непосредственно с кэшом
	}

	BrokerWithClient interface {
		Connect(ctx context.Context, brokers []string) error
		Close() error
	}
)
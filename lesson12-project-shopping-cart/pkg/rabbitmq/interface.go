package rabbitmq

import "context"

type RabbitMQService interface {
	Publish(ctx context.Context, queue string, message any) error
	Consume(ctx context.Context, queue string, handler func([]byte) error) error
	Close() error
}

package nredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Config configuration
type Config struct {
	Addr string
	Pass string
	DB   int
}

// Client Redis wrapper to handle pub/sub calls
type Client interface {
	Ping(ctx context.Context) (string, error)
	// Publish into Pub/Sub
	Publish(ctx context.Context, channel string, message []byte) error
	// Subscribe into Pub/Sub
	Subscribe(ctx context.Context, channel string, notifyTo chan string)
}

type client struct {
	C *redis.Client
}

// InitClient get a redis wrapper instance
func InitClient(config Config) (Client, error) {
	addr := "localhost:6379"
	if config.Addr != "" {
		addr = config.Addr
	}

	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Pass,
		DB:       10,
	})

	cli := client{
		C: c,
	}
	return &cli, nil
}

// Ping ping redis server
func (w *client) Ping(ctx context.Context) (string, error) {
	return w.C.Ping(ctx).Result()
}

// Publish message into a channel
func (w *client) Publish(ctx context.Context, channel string, message []byte) error {
	return w.C.Publish(ctx, channel, message).Err()
}

// Subscribe subscribe to listen messages from a channel
func (w *client) Subscribe(ctx context.Context, channel string, notifyTo chan string) {
	sub := w.C.Subscribe(ctx, channel)
	ch := sub.Channel()
	for msg := range ch {
		notifyTo <- msg.Payload
	}
}

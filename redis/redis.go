package redis

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

var ctx = context.Background()

// Client Redis wrapper to handle pub/sub calls
type Client interface {
	CheckHealth() (string, error)
	// Publish into Pub/Sub
	Publish(channel string, message []byte) error
	// Subscribe into Pub/Sub
	Subscribe(channel string, notifyTo chan string)
}

type client struct {
	KeyPrefix string
	rc        *redis.Client
}

// New get a redis wrapper instance
func New(config Config) (Client, error) {
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
		rc: c,
	}
	return &cli, nil
}

// CheckHealth ping redis server
func (init *client) CheckHealth() (string, error) {
	return init.rc.Ping(ctx).Result()
}

// Publish message into a channel
func (init *client) Publish(channel string, message []byte) error {
	return init.rc.Publish(ctx, channel, message).Err()
}

// Subscribe subscribe to listen messages from a channel
func (init *client) Subscribe(channel string, notifyTo chan string) {
	sub := init.rc.Subscribe(ctx, channel)
	ch := sub.Channel()
	for msg := range ch {
		notifyTo <- msg.Payload
	}
}

// WipeDB wipes the db.
func (init *client) WipeDB() error {
	err := init.rc.FlushDB(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

// Close terminates any storage connections gracefully.
func (init *client) Close() error {
	return init.rc.Close()
}

func (init *client) GetRedisPrefixedKey(key string) string {
	if init.KeyPrefix != "" {
		return init.KeyPrefix + ":" + key
	}
	return key
}

package libredis

import (
	"context"
	"encoding/json"
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
	Publish(channel string, message interface{}) error
	// Subscribe into Pub/Sub
	Subscribe(channel string, notifyTo chan string) error
	// Encode json encode
	Encode(model interface{}) ([]byte, error)
	// Decode json decode
	Decode(msg string, model interface{}) error

	WipeDB() error

	Close() error

	GetPrefixedKey(key string) string
}

type client struct {
	KeyPrefix string
	rc        *redis.Client
}

// New get redis wrapper instance
func New(config *Config) (Client, error) {
	addr := "localhost:6379"
	if config.Addr != "" {
		addr = config.Addr
	}

	_, err := redis.ParseURL(addr)
	if err != nil {
		panic(err)
	}

	c := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     10,
		MinIdleConns: 10,
		DB:           0,
	})
	cli := client{
		rc: c,
	}
	return &cli, nil
}

// CheckHealth ping server
func (init *client) CheckHealth() (string, error) {
	return init.rc.Ping(ctx).Result()
}

// Publish message into a channel
func (init *client) Publish(channel string, message interface{}) error {
	enc, err := init.Encode(message)
	if err != nil {
		return err
	}
	return init.rc.Publish(ctx, channel, enc).Err()
}

// Subscribe subscribe to listen messages from a channel
func (init *client) Subscribe(channel string, notifyTo chan string) error {
	sub := init.rc.Subscribe(ctx, channel)
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}
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

func (init *client) GetPrefixedKey(key string) string {
	if init.KeyPrefix != "" {
		return init.KeyPrefix + ":" + key
	}
	return key
}

func (init *client) Encode(model interface{}) ([]byte, error) {
	return json.Marshal(model)

}

func (init *client) Decode(msg string, model interface{}) error {
	if err := json.Unmarshal([]byte(msg), model); err != nil {
		return err
	} else {
		return nil
	}

}

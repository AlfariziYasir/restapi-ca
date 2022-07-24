package redis

import (
	"context"
	"fmt"
	"restapi/internal/config"

	redis "github.com/go-redis/redis/v8"
)

type Client interface {
	Conn() *redis.Client
	Close() error
}

func NewClientRDS() (Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Cfg().RedisHost, config.Cfg().RedisPort),
		DB:   config.Cfg().RedisDB,
	})
	err := db.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func NewClient() (Client, error) {
	return NewClientRDS()
}

type client struct {
	db *redis.Client
}

func (c *client) Conn() *redis.Client { return c.db }
func (c *client) Close() error        { return c.db.Close() }

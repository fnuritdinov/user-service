package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	//memoryCache "github.com/patrickmn/go-cache"

	"github.com/redis/go-redis/v9"
)

type CacheMemory struct {
	Name     string
	Email    string
	Password string
	Age      int32
	Phone    string
	Role     string
}
type ICache interface {
	Save(ctx context.Context, key string, data any, duration time.Duration) error
	Get(ctx context.Context, key string, data any) error
}

type cache struct {
	client *redis.Client
}

var ErrNotFound = errors.New("not found")

func New(ctx context.Context, address string) (ICache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	Ping := client.Ping(ctx)

	log.Println("client.Ping", Ping)

	return &cache{
		client: client,
	}, nil
}

func (c *cache) Save(ctx context.Context, key string, data any, duration time.Duration) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal %w", err)
	}

	err = c.client.Set(ctx, key, dataBytes, duration).Err()
	if err != nil {
		return fmt.Errorf(" c.client.Set %w", err)
	}

	return nil
}

func (c *cache) Get(ctx context.Context, key string, data any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return fmt.Errorf("c.client.Get %w", err)
	}

	err = json.Unmarshal([]byte(val), data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal %w", err)
	}

	return nil
}

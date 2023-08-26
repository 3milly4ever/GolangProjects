package caching

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration time.Duration) error
	Delete(key string) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client}
}

func (rc *RedisCache) Get(key string) ([]byte, error) {
	//Background is what ctx always gets initialized with
	ctx := context.Background()
	//we can use the ctx here to:
	//1.specify timeouts and cancellation signals for operations.
	//if our cache operation isn't supposed to take longer than a certain duration,
	//we can specify a timeout with context that can prevent the operation from
	//blocking indefinitely.this is a useful way to avoid waiting for a cache operation
	//that is slow or unresponsive.
	//2.the context can carry values across API boundaries and goroutines.
	//even though it's not utilized here, it's possible to add values to the context
	//that can be used downstream. values such as request IDs, user tokens, or other context-specific data.
	//3.context can signal a graceful shutdown. by canceling the context we indicate to our application that
	//it's time to wrap up any ongoing work early and make a clean exit.
	cachedValue, err := rc.getWithCtx(ctx, key)
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(cachedValue), nil
}

func (rc *RedisCache) getWithCtx(ctx context.Context, key string) (string, error) {
	cmd := rc.client.Get(key) //we use the Redis pack get command to get the key from our cache
	//the value we get with the key will be stored in cachedValue
	cachedValue, err := cmd.Result()
	if err != nil {
		//check if the error is within the reddis Get function which indicates that the key is not found in the cache
		if err == redis.Nil {
			return "", nil //return an empty string and no error for cache miss
		}
		return "", err //return actual error for other error cases
	}
	return cachedValue, nil //returns the cached value
}

func (rc *RedisCache) Set(key string, value []byte, expiration time.Duration) error {
	// ctx := context.Background()
	// args := []interface{}{value}
	return rc.client.Set(key, value, expiration).Err()
}

func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(key).Err()
}

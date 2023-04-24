package redis

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
	"time"
)

const (
	maxWait     = 3 * time.Second
	minInterval = 1 * time.Millisecond
	trySleep    = 100 * time.Millisecond
)

type Client struct {
	redis2.Client
}

func (c *Client) Lock(ctx context.Context, key string, expire time.Duration) (*Lock, error) {
	if expire < minInterval {
		expire = minInterval
	}
	ok, err := c.Set(ctx, key, "1", expire).Result()
	if ok == "OK" {
		return &Lock{
			key: key,
			c:   c,
		}, nil
	}
	return nil, err
}

func (c *Client) LockWait(ctx context.Context, key string, waitMax, expire time.Duration) (*Lock, error) {
	if waitMax > maxWait {
		waitMax = maxWait
	}
	endAt := time.Now().Add(waitMax)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			lock, err := c.Lock(ctx, key, expire)
			if lock != nil {
				return lock, err
			}
			startAt := time.Now()
			if startAt.Add(trySleep).After(endAt) {
				return nil, nil
			}
			time.Sleep(trySleep)
		}
	}
}

type Lock struct {
	key string
	c   *Client
}

func (r *Lock) Unlock(ctx context.Context) error {
	_, err := r.c.Del(ctx, r.key).Result()
	return err
}

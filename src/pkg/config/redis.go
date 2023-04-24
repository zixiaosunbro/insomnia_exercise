package config

import (
	redis2 "github.com/go-redis/redis/v8"
	"insomnia/src/pkg/middleware_custom/redis"
)

type RedisManager struct {
	redisPools map[string]*redis.Client
}

func NewRedisManager(redisConf map[string]*RedisConf) (*RedisManager, error) {
	pools := make(map[string]*redis.Client, len(redisConf))
	for name, redisCfg := range redisConf {
		option := &redis2.Options{
			DialTimeout:  redisCfg.ConnectTimeout,
			ReadTimeout:  redisCfg.ReadTimeout,
			WriteTimeout: redisCfg.WriteTimeout,
			PoolSize:     redisCfg.PoolMaxActive,
			IdleTimeout:  redisCfg.PoolIdleTimeout,
			MaxRetries:   redisCfg.MaxRetries,
			DB:           redisCfg.DB,
			Password:     redisCfg.Auth,
			Addr:         redisCfg.Uri,
		}
		client := redis2.NewClient(option)
		pools[name] = &redis.Client{Client: *client}

	}
	return &RedisManager{
		redisPools: pools,
	}, nil
}

func (rm *RedisManager) GetRedisPool(name string) *redis.Client {
	pool, exist := rm.redisPools[name]
	if !exist {
		panic("redis pool not exist")
	}
	return pool
}

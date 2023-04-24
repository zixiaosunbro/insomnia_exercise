package middleware

import (
	"insomnia/src/pkg/config"
	"insomnia/src/pkg/utils"
	"sync"
)

var (
	middleware = &struct {
		DBManager *config.Manager
		RedisM    *config.RedisManager
	}{}
	once = sync.Once{}
)

func Init() {
	once.Do(func() {
		cfg := config.GetConfig()
		// plugin like a switch, if mysql no longer used, set it to false. service will not init mysql
		plugins := cfg.Plugins

		// init mysql
		if plugins.DB {
			dbManager, err := config.NewManager(cfg.DB)
			utils.Must(err)
			middleware.DBManager = dbManager
		}
		// init redis
		if plugins.Redis {
			redisManager, err := config.NewRedisManager(cfg.Redis)
			utils.Must(err)
			middleware.RedisM = redisManager
		}

	})
}

func GetDBManager() *config.Manager {
	if middleware.DBManager == nil {
		panic("db not init")
	}
	return middleware.DBManager
}

func GetRedisManager() *config.RedisManager {
	if middleware.RedisM == nil {
		panic("redis not init")
	}
	return middleware.RedisM
}

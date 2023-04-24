package rule

import (
	"gorm.io/gorm"
	"insomnia/src/pkg/config/middleware"
	"insomnia/src/pkg/middleware_custom/redis"
	"sync"
)

var (
	once     sync.Once
	DBMaster *gorm.DB
	DBSlave  *gorm.DB
	RedisCli *redis.Client
)

func ResourceInit() {
	once.Do(func() {
		DBMaster = middleware.GetDBManager().GetMasterDB("insomnia")
		DBSlave = middleware.GetDBManager().GetSlaveDB("insomnia")
		RedisCli = middleware.GetRedisManager().GetRedisPool("insomnia")
	})
}

package utils

import (
	"github.com/go-redis/redis/v8"
	"strings"
)

func MysqlReturnNil(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "no rows in result set") {
		return true
	}
	return false
}

func RedisReturnNil(err error) bool {
	return err == redis.Nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

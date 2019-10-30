package dbservice

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

type RedisDB struct {
	pool *redis.Pool // Redis 连接池
	DBValid bool     // Redis是否正常
}


func (redisDB *RedisDB) initPool() {
	if redisDB.pool == nil {
		addr := viper.GetString("Redis.Addr")
		index := viper.GetString("Redis.Index")
		rawUrl := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("Redis.Password")

		maxIdle := viper.GetInt("Redis.MaxIdle")
		idleTimeout := viper.GetInt("Reids.IdleTimeout")
		maxActive := viper.GetInt("Redis.MaxActive")

		redisDB.pool = &redis.Pool{
			MaxIdle: maxIdle,
			MaxActive: maxActive,
			Wait: false,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawUrl, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawUrl)
			},
		}
	}
}

func (redisDB *RedisDB) IsDBRedisValid() bool {
	c := redisDB.Get()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false;
	}
	return true;
}

func (redisDB *RedisDB) checkHealth() {
	t := viper.GetInt("Redis.RedisHealthCheckTimer")
	if t == 0 {
		t = 1
	}

	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			redisDB.DBValid = redisDB.IsDBRedisValid()
		}
	}
}

func (redisDB *RedisDB) Get() redis.Conn {
	if redisDB.pool == nil {
		redisDB.initPool()
		go redisDB.checkHealth()
	}
	return redisDB.pool.Get()
}
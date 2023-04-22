package redis

import (
	"github.com/go-redis/redis"
	"github.com/k4zb3k/project/config"
	"github.com/k4zb3k/project/pkg/logger"
	"net"
)

func InitRedis(cfg config.CacheConnConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logger.Error.Fatalln(err)
		return nil, err
	}
	return client, nil
}

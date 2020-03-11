package redis

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

// Init the redis client
func Init() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return nil, errors.New("no redis host was provided")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPass := os.Getenv("REDIS_PASS")

	redisClient := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), Password: redisPass, TLSConfig: &tls.Config{ServerName: redisHost}, DB: 0})
	_, err := redisClient.Ping().Result()
	if err != nil && os.Getenv("REDIS_ALLOW_INSECURE") == "true" {
		redisClient = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), Password: redisPass, DB: 0})
		_, err = redisClient.Ping().Result()
	}
	if err != nil {
		return nil, fmt.Errorf("unable to ping redis: %s", err.Error())
	}
	return redisClient, nil
}

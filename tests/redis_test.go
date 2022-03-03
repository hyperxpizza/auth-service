package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

//go test -v ./tests --run TestRedisConnection --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestRedisConnection(t *testing.T) {
	flag.Parse()

	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:   int(cfg.Redis.DB),
	})

	stat, err := rdc.Ping(ctx).Result()
	assert.NoError(t, err)

	assert.Equal(t, stat, "PONG")
}
